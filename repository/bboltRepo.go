package repository

import (
	"time"

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
	return nil
}

func (*bboltRepo) GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, error) {
	var worklogs []*model.Work

	return worklogs, nil
}

func (*bboltRepo) GetByID(ID string, filter *model.Work) (*model.Work, error) {
}

// Internal wrapped function to ensure all useages are aligned
func open() (*bolt.DB, error) {
	db, err := bolt.Open(filePath, 0750, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	return db, nil
}
