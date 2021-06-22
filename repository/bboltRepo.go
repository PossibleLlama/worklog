package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"time"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/model"

	bolt "go.etcd.io/bbolt"
)

var filePath string

const worklogBucket = "worklogs"

type bboltRepo struct{}

// NewBBoltRepo initializes the repo with the given filepath
func NewBBoltRepo(path string) WorklogRepository {
	filePath = path
	return &bboltRepo{}
}

func (*bboltRepo) Configure(cfg *model.Config) error {
	return nil
}

func (*bboltRepo) Save(wl *model.Work) error {
	db, openErr := open()
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	updateErr := db.Update(func(tx *bolt.Tx) error {
		b, bucketErr := tx.CreateBucketIfNotExists([]byte(worklogBucket))
		if bucketErr != nil {
			return bucketErr
		}
		created, marshalErr := json.Marshal(wl)
		if marshalErr != nil {
			return marshalErr
		}
		return b.Put([]byte(wl.ID), created)
	})
	return updateErr
}

func (*bboltRepo) GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, error) {
	var worklogs []*model.Work
	db, openErr := openReadOnly()
	if openErr != nil {
		return nil, openErr
	}
	defer db.Close()

	viewErr := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(worklogBucket))
		if b == nil {
			return errors.New(e.RepoGetFilesRead)
		}
		c := b.Cursor()
		min, minErr := startDate.MarshalBinary()
		if minErr != nil {
			return minErr
		}
		max, maxErr := endDate.MarshalBinary()
		if maxErr != nil {
			return maxErr
		}

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			var found model.Work
			if marshalErr := json.Unmarshal(v, &found); marshalErr != nil {
				return marshalErr
			}

			if workMatchesFilter(filter, &found) {
				worklogs = append(worklogs, &found)
			}
		}
		return nil
	})

	return worklogs, viewErr
}

func (*bboltRepo) GetByID(ID string, filter *model.Work) (*model.Work, error) {
	var foundWl *model.Work
	db, openErr := openReadOnly()
	if openErr != nil {
		return nil, openErr
	}
	defer db.Close()

	viewErr := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(worklogBucket))
		found := b.Get([]byte(ID))
		return json.Unmarshal(found, foundWl)
	})
	return foundWl, viewErr
}

// Internal wrapped function to ensure all useages are aligned
func open() (*bolt.DB, error) {
	return bolt.Open(filePath, 0750, &bolt.Options{
		Timeout: 1 * time.Second,
	})
}

// Internal wrapped function to ensure all useages are aligned
func openReadOnly() (*bolt.DB, error) {
	return bolt.Open(filePath, 0750, &bolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: true,
	})
}
