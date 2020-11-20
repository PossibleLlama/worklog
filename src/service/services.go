package service

import (
	"net/http"

	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
)

// WorklogService defines what a service for
// worklogs should be capable of doing
type WorklogService interface {
	CreateWorklog(wl *model.Work) (int, error)
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

func (*service) CreateWorklog(wl *model.Work) (int, error) {
	if err := repo.Save(wl); err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusCreated, nil
}
