package repository

import (
	"time"

	"github.com/PossibleLlama/worklog/model"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	bolt "go.etcd.io/bbolt"
)

var filePath string

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
	db, openErr := openReadWrite()
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	return db.Save(wl)
}

func (*bboltRepo) GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, error) {
	var worklogs []*model.Work
	db, openErr := openReadOnly()
	if openErr != nil {
		return nil, openErr
	}
	defer db.Close()

	sel := q.And(
		q.Gte("When", startDate),
		q.Lt("When", endDate),
		filterQuery(filter),
	)
	viewErr := db.Select(sel).OrderBy("When").Find(&worklogs)

	if viewErr == storm.ErrNotFound {
		return worklogs, nil
			}
	return worklogs, viewErr
}

func (*bboltRepo) GetByID(ID string, filter *model.Work) (*model.Work, error) {
	var foundWl model.Work
	db, openErr := openReadOnly()
	if openErr != nil {
		return nil, openErr
	}
	defer db.Close()

	sel := q.And(
		q.Re("ID", ID),
		filterQuery(filter),
	)
	viewErr := db.Select(sel).OrderBy("Revision").Limit(1).First(&foundWl)

	if viewErr == storm.ErrNotFound {
		return &foundWl, nil
	}
	return &foundWl, viewErr
}

// Internal wrapped function to ensure all useages are aligned
func openReadWrite() (*storm.DB, error) {
	return storm.Open(filePath, storm.BoltOptions(0750, &bolt.Options{
		Timeout: 1 * time.Second,
	}))
}

// Internal wrapped function to ensure all useages are aligned
func openReadOnly() (*storm.DB, error) {
	return storm.Open(filePath, storm.BoltOptions(0750, &bolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: true,
	}))
}

func filterQuery(f *model.Work) q.Matcher {
	sel := q.And(
		q.Re("Title", f.Title),
		q.Re("Description", f.Description),
		q.Re("Author", f.Author),
	)
	return sel
}
