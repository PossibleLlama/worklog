package repository

import (
	"errors"
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

type bboltRepo struct{}

// NewBBoltRepo initializes the repo with the given filepath
func NewBBoltRepo(path string) WorklogRepository {
	filePath = path
	return &bboltRepo{}
}

func (*bboltRepo) Save(wl *model.Work) error {
	db, openErr := openReadWrite()
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	helpers.LogDebug("Saving file...")
	if err := db.Save(wl); err != nil {
		return err
	}

	helpers.LogDebug("Saved file")
	return nil
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
	var foundWls []*model.Work
	db, openErr := openReadOnly()
	if openErr != nil {
		return nil, openErr
	}
	defer db.Close()

	sel := q.And(
		q.Re("ID", helpers.RegexCaseInsensitive+ID),
		filterQuery(filter),
	)
	viewErr := db.Select(sel).OrderBy("Revision").Find(&foundWls)

	if viewErr == storm.ErrNotFound {
		return nil, nil
	} else if viewErr == nil && len(foundWls) > 1 {
		return nil, errors.New(e.RepoGetSingleFileAmbiguous)
	}
	if !filterByTags(filter, foundWls[0]) {
		return nil, viewErr
	}
	return foundWls[0], viewErr
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
		q.Re("Title", helpers.RegexCaseInsensitive+f.Title),
		q.Re("Description", helpers.RegexCaseInsensitive+f.Description),
		q.Re("Author", helpers.RegexCaseInsensitive+f.Author),
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
