package repository

import (
	"fmt"
	"os"
	"strings"
	"time"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	bolt "go.etcd.io/bbolt"
)

var filePath string

const regexCaseInsesitive = "(?i)"

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
	var foundWls, filteredWls []*model.Work
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
	viewErr := db.Select(sel).OrderBy("When").Find(&foundWls)

	for _, el := range foundWls {
		if filterByTags(filter, el) {
			filteredWls = append(filteredWls, el)
		}
	}

	if viewErr == storm.ErrNotFound {
		return []*model.Work{}, nil
	}
	return filteredWls, viewErr
}

func (*bboltRepo) GetByID(ID string, filter *model.Work) (*model.Work, error) {
	var foundWl *model.Work
	db, openErr := openReadOnly()
	if openErr != nil {
		return nil, openErr
	}
	defer db.Close()

	sel := q.And(
		q.Re("ID", regexCaseInsesitive+ID),
		filterQuery(filter),
	)
	viewErr := db.Select(sel).OrderBy("Revision").First(&foundWl)

	if viewErr == storm.ErrNotFound {
		return foundWl, nil
	}
	if !filterByTags(filter, foundWl) {
		return nil, viewErr
	}
	return foundWl, viewErr
}

// Internal wrapped function to ensure all useages are aligned
func openReadWrite() (*storm.DB, error) {
	return storm.Open(filePath, storm.BoltOptions(0750, &bolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: false,
	}))
}

// Internal wrapped function to ensure all useages are aligned
func openReadOnly() (*storm.DB, error) {
	if _, err := os.Stat(filePath); err == nil {
		return storm.Open(filePath, storm.BoltOptions(0750, &bolt.Options{
			Timeout:  1 * time.Second,
			ReadOnly: true,
		}))
	}
	return nil, fmt.Errorf(e.RepoGetFilesRead)
}

func filterQuery(f *model.Work) q.Matcher {
	sel := q.And(
		q.Re("Title", regexCaseInsesitive+f.Title),
		q.Re("Description", regexCaseInsesitive+f.Description),
		q.Re("Author", regexCaseInsesitive+f.Author),
	)
	return sel
}

// Check if the filters tags exist in the provided work.
// Returns true if it matches, or false if not.
func filterByTags(f *model.Work, w *model.Work) bool {
	if f == nil || w == nil {
		return false
	}
	for _, fTags := range f.Tags {
		if fTags != "" && !helpers.AInB(fTags, strings.Join(w.Tags, " ")) {
			return false
		}
	}
	return true
}
