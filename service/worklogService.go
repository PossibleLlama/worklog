package service

import (
	"net/http"
	"sort"
	"time"

	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
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

func (*service) Configure(cfg *model.Config) error {
	return repo.Configure(cfg)
}

func (*service) CreateWorklog(wl *model.Work) (int, error) {
	if err := repo.Save(wl); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func (*service) GetWorklogsBetween(start, end time.Time, filter *model.Work) ([]*model.Work, int, error) {
	worklogs := make(model.WorkList, 0)
	worklogs, err := repo.GetAllBetweenDates(start, end, filter)
	if err != nil {
		return worklogs, http.StatusInternalServerError, err
	}
	worklogs = removeOldRevisions(worklogs)
	if len(worklogs) == 0 {
		return worklogs, http.StatusNotFound, err
	}
	sort.Sort(worklogs)
	return worklogs, http.StatusOK, nil
}

func removeOldRevisions(wls []*model.Work) []*model.Work {
	deDuplicated := []*model.Work{}
	uniqueIDWls := make(map[string][]*model.Work)
	for _, element := range wls {
		uniqueIDWls[element.ID] = append(uniqueIDWls[element.ID], element)
	}
	for _, wls := range uniqueIDWls {
		highestRevision := -1
		for _, element := range wls {
			if element.Revision > highestRevision {
				highestRevision = element.Revision
			}
		}
		for _, element := range wls {
			if element.Revision == highestRevision {
				deDuplicated = append(deDuplicated, element)
				break
			}
		}
	}
	return deDuplicated
}

func (*service) GetWorklogsByID(filter *model.Work, ids ...string) ([]*model.Work, int, error) {
	worklogs := make(model.WorkList, 0)
	for _, ID := range ids {
		// Implement using goroutines and channels
		wl, err := repo.GetByID(ID, filter)
		if err != nil {
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
