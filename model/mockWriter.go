package model

import (
	"github.com/stretchr/testify/mock"
)

// mockWriter mock to satisfy io.Writer
type mockWriter struct {
	mock.Mock
}

// Write to satisfy io.Writer
func (m *mockWriter) Write(p []byte) (int, error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}
