package repository

import (
	"time"

	"github.com/PossibleLlama/worklog/model"
)

// WorklogRepository defines what a repository
// for worklogs should be capable of doing
type WorklogRepository interface {
	Init() error
	Save(wl *model.Work) error

	GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, error)
	GetByID(id string, filter *model.Work) (*model.Work, error)
	GetAll() ([]*model.Work, error)
}

// ConfigRepository defines what a configuration
// store should be capable of doing
type ConfigRepository interface {
	SaveConfig(cfg *model.Config) error
}
