package repository

import (
	"github.com/PossibleLlama/worklog/model"
)

// WorklogRepository defines what a repository
// for worklogs should be capable of doing
type WorklogRepository interface {
	Save(wl *model.Work) error
}
