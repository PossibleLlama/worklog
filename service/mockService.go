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

// CreateWorklog WorklogService method for testing
func (m *MockService) CreateWorklog(wl *model.Work) (int, error) {
	args := m.Called(wl)
	return args.Int(0), args.Error(1)
}

// EditWorklog WorklogService method for testing
func (m *MockService) EditWorklog(id string, newWl *model.Work) (*model.Work, int, error) {
	args := m.Called(id, newWl)
	return nil, args.Int(0), args.Error(1)
}

// GetWorklogsBetween WorklogService method for testing
func (m *MockService) GetWorklogsBetween(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, int, error) {
	args := m.Called(startDate, endDate, filter)
	return args.Get(0).([]*model.Work), args.Int(1), args.Error(2)
}

// GetWorklogsByID WorklogService method for testing
func (m *MockService) GetWorklogsByID(filter *model.Work, ids ...string) ([]*model.Work, int, error) {
	args := m.Called(filter, ids)
	return args.Get(0).([]*model.Work), args.Int(1), args.Error(2)
}

// ExportTo WorklogService method for testing
func (m *MockService) ExportTo(path string) (int, error) {
	args := m.Called(path)
	return args.Int(0), args.Error(1)
}
