package repository

import (
	"time"

	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/mock"
)

// MockRepo type of WorklogRepository for testing
type MockRepo struct {
	mock.Mock
}

// SaveConfig WorklogRepository method for testing
func (m *MockRepo) SaveConfig(cfg *model.Config) error {
	args := m.Called(cfg)
	return args.Error(0)
}

func (m *MockRepo) Init() error {
	args := m.Called()
	return args.Error(0)
}

// Save WorklogRepository method for testing
func (m *MockRepo) Save(wl *model.Work) error {
	args := m.Called(wl)
	return args.Error(0)
}

// GetAllBetweenDates WorklogRepository method for testing
func (m *MockRepo) GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, error) {
	args := m.Called(startDate, endDate, filter)
	return args.Get(0).([]*model.Work), args.Error(1)
}

// GetByID WorklogRepository method for testing
func (m *MockRepo) GetByID(id string, filter *model.Work) (*model.Work, error) {
	args := m.Called(id, filter)
	return args.Get(0).(*model.Work), args.Error(1)
}

// GetAll WorklogRepository method for testing
func (m *MockRepo) GetAll() ([]*model.Work, error) {
	args := m.Called()
	return args.Get(0).([]*model.Work), args.Error(1)
}
