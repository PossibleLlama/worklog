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

// Configure WorklogRepository method for testing
func (m *MockRepo) Configure(cfg *model.Config) error {
	args := m.Called(cfg)
	return args.Error(0)
}

// Save WorklogRepository method for testing
func (m *MockRepo) Save(wl *model.Work) error {
	args := m.Called(wl)
	return args.Error(0)
}

// GetAllBetweenDates WorklogRepository method for testing
func (m *MockRepo) GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) (model.WorkList, error) {
	args := m.Called(startDate, endDate, filter)
	return args.Get(0).([]*model.Work), args.Error(1)
}
