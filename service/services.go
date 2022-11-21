package service

import (
	"time"

	"github.com/PossibleLlama/worklog/model"
)

// WorklogService defines what a service for worklogs should be capable of doing
type WorklogService interface {
	CreateWorklog(wl *model.Work) (int, error)
	EditWorklog(id string, newWl *model.Work) (*model.Work, int, error)
	GetWorklogsBetween(start, end time.Time, filter *model.Work) ([]*model.Work, int, error)
	GetWorklogsByID(filter *model.Work, ids ...string) ([]*model.Work, int, error)
}
