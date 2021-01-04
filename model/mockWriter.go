package model

import (
	"github.com/stretchr/testify/mock"
)

// MockWriter mock to satisfy io.Writer
type MockWriter struct {
	mock.Mock
}

// Write to satisfy io.Writer
func (m *MockWriter) Write(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}
