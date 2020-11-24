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
	Congfigure(author string, duration int) error
	CreateWorklog(wl *model.Work) (int, error)
	GetWorklogsSince(date time.Time) ([]*model.Work, int, error)
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

func (*service) Congfigure(author string, duration int) error {
	return repo.Configure(author, duration)
}

func (*service) CreateWorklog(wl *model.Work) (int, error) {
	if err := repo.Save(wl); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}

func (*service) GetWorklogsSince(date time.Time) ([]*model.Work, int, error) {
	worklogs, err := repo.GetAllSinceDate(date)
	if err != nil {
		return worklogs, http.StatusInternalServerError, err
	}
	return worklogs, http.StatusOK, nil
}
