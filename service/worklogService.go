package service

import (
	"net/http"
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
	worklogs, err := repo.GetAllBetweenDates(start, end, filter)
	worklogs = removeOldRevisions(worklogs)
	if err != nil {
		return worklogs, http.StatusInternalServerError, err
	}
	if len(worklogs) == 0 {
		return worklogs, http.StatusNotFound, err
	}
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
