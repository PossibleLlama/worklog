package service

import (
	"net/http"
	"sort"
	"strings"
	"time"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"

	"github.com/spf13/viper"
)

type service struct{}

var (
	repo repository.WorklogRepository
)

// NewWorklogService Generator for service while initialising the storage capability
func NewWorklogService(repository repository.WorklogRepository) WorklogService {
	repo = repository
	return &service{}
}

func (*service) CreateWorklog(wl *model.Work) (int, error) {
	if wl.Duration <= 0 {
		wl.Duration = viper.GetInt("default.duration")
	}
	wl.Title = helpers.Sanitize(strings.TrimSpace(wl.Title))
	wl.Description = helpers.Sanitize(strings.TrimSpace(wl.Description))
	wl.Author = helpers.Sanitize(strings.TrimSpace(wl.Author))

	cleanTags := make([]string, 0)
	for _, tag := range wl.Tags {
		if strings.TrimSpace(tag) != "" {
			cleanTags = append(cleanTags, helpers.Sanitize(strings.TrimSpace(tag)))
		}
	}
	wl.Tags = cleanTags
	if wl.Author == "" {
		wl.Author = helpers.Sanitize(viper.GetString("default.author"))
	}

	wl.Tags = helpers.DeduplicateString(wl.Tags)
	if wl.When.IsZero() {
		wl.When = wl.CreatedAt
	}
	if err := repo.Save(wl); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func (s *service) EditWorklog(id string, newWl *model.Work) (*model.Work, int, error) {
	newWl.Tags = helpers.DeduplicateString(newWl.Tags)
	wls, code, err := s.GetWorklogsByID(&model.Work{}, id)
	if err != nil {
		return nil, code, err
	}
	// The get returns 1 WL per ID, as we are only
	// passing one ID, there is only one possible WL
	if len(wls) == 0 {
		return nil, http.StatusNotFound, nil
	}

	wl := wls[0]
	wl.Update(*newWl)

	if err := repo.Save(wl); err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return wl, http.StatusOK, nil
}

func (*service) GetWorklogsBetween(start, end time.Time, filter *model.Work) ([]*model.Work, int, error) {
	var worklogs model.WorkList
	var err error

	if (end == time.Time{}) {
		end = time.Date(3000, time.January, 1, 0, 0, 0, 0, time.Now().Location())
	}

	worklogs, err = repo.GetAllBetweenDates(start, end, filter)
	if err != nil {
		return worklogs, http.StatusInternalServerError, err
	}
	worklogs = worklogs.RemoveOldRevisions()
	if len(worklogs) == 0 {
		return worklogs, http.StatusNotFound, err
	}
	sort.Sort(worklogs)
	return worklogs, http.StatusOK, nil
}

func (*service) GetWorklogsByID(filter *model.Work, ids ...string) ([]*model.Work, int, error) {
	worklogs := make(model.WorkList, 0)
	for _, ID := range ids {
		// Implement using goroutines and channels
		wl, err := repo.GetByID(ID, filter)
		if err != nil {
			if err.Error() == e.RepoGetSingleFileAmbiguous {
				return worklogs, http.StatusNotFound, nil
			}
			return worklogs, http.StatusInternalServerError, err
		}
		if wl != nil {
			worklogs = append(worklogs, wl)
		}
	}
	if len(worklogs) == 0 {
		return worklogs, http.StatusNotFound, nil
	}
	sort.Sort(worklogs)
	return worklogs, http.StatusOK, nil
}
