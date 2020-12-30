package service

import (
	"net/http"
	"time"

	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
)

// WorklogService defines what a service for
// worklogs should be capable of doing
type WorklogService interface {
	Congfigure(cfg *model.Config) error
	CreateWorklog(wl *model.Work) (int, error)
	GetWorklogsBetween(start, end time.Time, filter *model.Work) ([]*model.Work, int, error)
}

type service struct{}

var (
	repo repository.WorklogRepository
)

// NewWorklogService Generator for service while initialising the storage capability
func NewWorklogService(repository repository.WorklogRepository) WorklogService {
	repo = repository
	return &service{}
}

func (*service) Congfigure(cfg *model.Config) error {
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
	if err != nil {
		return worklogs, http.StatusInternalServerError, err
	}
	return worklogs, http.StatusOK, nil
}
