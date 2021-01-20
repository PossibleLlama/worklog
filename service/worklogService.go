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
	wls = reverse(wls)
	keys := make(map[string]bool)
	deDuplicated := []*model.Work{}
	for _, element := range wls {
		if _, value := keys[element.ID]; !value {
			keys[element.ID] = true
			deDuplicated = append(deDuplicated, element)
		}
	}
	return deDuplicated
}

func reverse(wls []*model.Work) []*model.Work {
	revList := []*model.Work{}
	for i := len(wls) - 1; i >= 0; i-- {
		revList = append(revList, wls[i])
	}
	return revList
}
