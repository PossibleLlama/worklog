package service

import (
	"time"

	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/mock"
)

// MockService type of WorklogService for testing
type MockService struct {
	mock.Mock
}

// Configure WorklogService method for testing
func (m *MockService) Configure(cfg *model.Config) error {
	args := m.Called(cfg)
	return args.Error(0)
}

// CreateWorklog WorklogService method for testing
func (m *MockService) CreateWorklog(wl *model.Work) (int, error) {
	args := m.Called(wl)
	return args.Int(0), args.Error(1)
}

// GetWorklogsBetween WorklogService method for testing
func (m *MockService) GetWorklogsBetween(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, int, error) {
	args := m.Called(startDate, endDate, filter)
	return args.Get(0).([]*model.Work), args.Int(0), args.Error(1)
}
