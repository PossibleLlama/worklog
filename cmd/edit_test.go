package cmd

import (
	"errors"
	"testing"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/service"

	"github.com/stretchr/testify/assert"
)

func setProvidedEditValues(title, description string, duration int, author, when, tags string) {
	editTitle = title
	editDescription = description
	editDuration = duration
	editAuthor = author
	editWhenString = when
	editTagsString = tags
}

func TestEditArgs(t *testing.T) {
	var tests = []struct {
		name string
		ids  []string
		err  error
	}{
		{
			name: "Single ID",
			ids:  []string{helpers.RandAlphabeticString(shortLength)},
			err:  nil,
		}, {
			name: "No args throws error",
			ids:  []string{},
			err:  errors.New(e.EditID),
		}, {
			name: "2 ids throws error",
			ids: []string{
				helpers.RandAlphabeticString(shortLength),
				helpers.RandAlphabeticString(shortLength)},
			err: errors.New(e.EditID),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedEditValues(
				helpers.RandAlphabeticString(shortLength),
				helpers.RandAlphabeticString(shortLength),
				shortLength,
				helpers.RandAlphabeticString(shortLength),
				"2020-01-01T13:24:50Z",
				helpers.RandAlphabeticString(shortLength))

			retErr := editArgs(testItem.ids)

			if testItem.err == nil {
				assert.Nil(t, retErr)
			} else {
				assert.EqualError(t, retErr, testItem.err.Error())
			}
		})
	}
}

func TestEditRun(t *testing.T) {
	var tests = []struct {
		name     string
		author   string
		duration int
		format   string
		expErr   error
	}{
		{
			name:     "Sends model to service",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: shortLength,
			format:   "yaml",
			expErr:   nil,
		}, {
			name:     "Error propogated",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: shortLength,
			format:   "yaml",
			expErr:   errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		cfg := model.NewConfig(testItem.author, testItem.format, testItem.duration)

		mockService := new(service.MockService)
		mockService.On("Configure", cfg).Return(testItem.expErr)
		wlService = mockService

		t.Run(testItem.name, func(t *testing.T) {
			setProvidedConfigureValues(testItem.author, testItem.format, testItem.duration)

			actualErr := configRun()

			mockService.AssertExpectations(t)
			mockService.AssertCalled(t,
				"Configure",
				cfg)
			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
