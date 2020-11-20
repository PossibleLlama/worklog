package service

import (
	"github.com/PossibleLlama/worklog/model"
)

// WorklogService defines what a service for
// worklogs should be capable of doing
type WorklogService interface {
	CreateWorklog(wl *model.Work) error
}

type service struct{}

var (
// repo   repository.WorklogRepository
)

// NewWorklogService Generator for service while initialising the storage capability
func NewWorklogService() WorklogService {
	// repo = repository
	return &service{}
}

func (*service) CreateWorklog(wl *model.Work) error {
	return nil
}
