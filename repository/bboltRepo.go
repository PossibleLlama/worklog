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

func (*bboltRepo) Init() error {
	var foundWls []*model.Work

	db, openErr := openReadWrite()
	if openErr != nil {
		helpers.LogError(fmt.Sprintf("Error opening file: %s", openErr.Error()), "read db error - bolt")
		return openErr
	}
	defer func() {
		_ = db.Close()
	}()

	viewErr := db.Select(
		q.Eq("WhenQueryEpoch", 0)).Find(&foundWls)
	if viewErr != nil && viewErr != storm.ErrNotFound {
		return errors.New("failed to get from db")
	}

	helpers.LogDebug(fmt.Sprintf("found %d items without epoch", len(foundWls)), "update db error - bolt")
	for _, el := range foundWls {
		helpers.LogDebug(
			fmt.Sprintf("ID %s, with when %s with epoch %d",
				el.ID, helpers.TimeFormat(el.When), el.WhenQueryEpoch),
			"update db - bolt")
		err := db.UpdateField(&model.Work{ID: el.ID}, "WhenQueryEpoch", el.When.Unix())
		if err != nil {
			helpers.LogError(
				fmt.Sprintf("failed to update element ID: %s's when epoch. error: %s", el.ID, err.Error()),
				"update db error - bolt")
		}
	}
	if err := db.ReIndex(&model.Work{}); err != nil {
		helpers.LogError("failed to reindex database", "update db - bolt")
	}

	return nil
}

func (*bboltRepo) Save(wl *model.Work) error {
	db, openErr := openReadWrite()
	if openErr != nil {
		helpers.LogError(fmt.Sprintf("Error opening file: %s", openErr.Error()), "read db error - bolt")
		return openErr
	}
	defer func() {
		_ = db.Close()
	}()

	helpers.LogDebug("Saving file...", "save model - bolt")
	if err := db.Save(wl); err != nil {
		helpers.LogError(fmt.Sprintf("Error closing file: %s", err.Error()), "save model error - bolt")
		return err
	}

	helpers.LogDebug("Saved file", "save model successful - bolt")
	return nil
}

func (*bboltRepo) GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, error) {
	var foundWls, filteredWls []*model.Work
	db, openErr := openReadOnly()
	if openErr != nil {
		helpers.LogError(fmt.Sprintf("Error opening file: %s", openErr.Error()), "read db error - bolt")
		return nil, openErr
	}
	defer func() {
		_ = db.Close()
	}()

	sel := q.And(
		q.Gte("WhenQueryEpoch", startDate.Unix()),
		q.Lt("WhenQueryEpoch", endDate.Unix()),
		filterQuery(filter),
	)
	viewErr := db.Select(sel).OrderBy("WhenQueryEpoch").Find(&foundWls)
	if viewErr != nil && viewErr != storm.ErrNotFound {
		return nil, errors.New("failed to get from db between dates")
	}

	for _, el := range foundWls {
		if filterByTags(filter, el) {
			el.Sanitize()
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
		helpers.LogError(fmt.Sprintf("Error opening file: %s", openErr.Error()), "read db error - bolt")
		return nil, openErr
	}
	defer func() {
		_ = db.Close()
	}()

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
	foundWls[0].Sanitize()
	return foundWls[0], viewErr
}

func (*bboltRepo) GetAll() ([]*model.Work, error) {
	var all []*model.Work
	db, openErr := openReadOnly()
	if openErr != nil {
		helpers.LogError(fmt.Sprintf("Error opening file: %s", openErr.Error()), "read db error - bolt")
		return nil, openErr
	}
	defer func() {
		_ = db.Close()
	}()

	err := db.All(&all)
	return all, err
}

// Internal wrapped function to ensure all usages are aligned
func openReadWrite() (*storm.DB, error) {
	return storm.Open(filePath, storm.BoltOptions(0750, &bolt.Options{
		Timeout:  1 * time.Second,
		ReadOnly: false,
	}))
}

// Internal wrapped function to ensure all usages are aligned
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
