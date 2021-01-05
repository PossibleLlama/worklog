package cmd

import (
	"errors"
	"testing"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/service"
	"github.com/stretchr/testify/assert"
)

const (
	shortLength = 30
	longLength  = 256
)

func setProvidedValues(author, format string, duration int) {
	configProvidedAuthor = author
	configProvidedDuration = duration
	configProvidedFormat = format
}

func TestConfigArgs(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{
			name: "Variables take defaults",
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedValues(helpers.RandString(shortLength), helpers.RandString(shortLength), shortLength)

			configArgs()

			assert.Equal(t, configDefaultAuthor, configProvidedAuthor)
			assert.Equal(t, configDefaultDuration, configProvidedDuration)
			assert.Equal(t, configDefaultFormat, configProvidedFormat)
		})
	}
}

func TestConfigRun(t *testing.T) {
	var tests = []struct {
		name     string
		author   string
		duration int
		format   string
		expErr   error
	}{
		{
			name:     "Sends model to service",
			author:   helpers.RandString(shortLength),
			duration: shortLength,
			format:   "yaml",
			expErr:   nil,
		}, {
			name:     "Error propogated",
			author:   helpers.RandString(shortLength),
			duration: shortLength,
			format:   "yaml",
			expErr:   errors.New(helpers.RandString(shortLength)),
		},
	}

	for _, testItem := range tests {
		cfg := model.NewConfig(testItem.author, testItem.format, testItem.duration)

		mockService := new(service.MockService)
		mockService.On("Congfigure", cfg).Return(testItem.expErr)
		wlService = mockService

		t.Run(testItem.name, func(t *testing.T) {
			setProvidedValues(testItem.author, testItem.format, testItem.duration)

			actualErr := configRun()

			mockService.AssertExpectations(t)
			mockService.AssertCalled(t,
				"Congfigure",
				cfg)
			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
