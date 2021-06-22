package repository

import (
	"time"

	"github.com/PossibleLlama/worklog/model"

	bolt "go.etcd.io/bbolt"
)

type bboltRepo struct{}

func NewBBoltRepo() WorklogRepository {
	return &bboltRepo{}
}

func (*bboltRepo) Configure(cfg *model.Config) error {
	return nil
}

func (*bboltRepo) Save(wl *model.Work) error {
	return nil
}

func (*bboltRepo) GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, error) {
	var worklogs []*model.Work

	return worklogs, nil
}

func (*bboltRepo) GetByID(ID string, filter *model.Work) (*model.Work, error) {
	return nil, nil
}
