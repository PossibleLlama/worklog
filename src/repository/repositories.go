package repository

import (
	"time"

	"github.com/PossibleLlama/worklog/model"
)

// WorklogRepository defines what a repository
// for worklogs should be capable of doing
type WorklogRepository interface {
	Configure(cfg *model.Config) error
	Save(wl *model.Work) error

	GetAllSinceDate(startDate time.Time) ([]*model.Work, error)
}
