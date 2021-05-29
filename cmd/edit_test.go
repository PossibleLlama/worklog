package cmd

import (
	"errors"
	"strings"
	"testing"
	"time"

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
	id := helpers.RandAlphabeticString(shortLength)
	now := time.Date(2020, time.January, 30, 23, 59, 0, 0, time.UTC)

	var tests = []struct {
		name     string
		ids      []string
		provided *model.Work
		expected *model.Work
		err      error
	}{
		{
			name:     "Single ID",
			ids:      []string{id},
			provided: nil,
			expected: &model.Work{
				Title:       "",
				Description: "",
				Duration:    0,
				When:        time.Time{},
				Tags:        []string{},
			},
			err: nil,
		}, {
			name:     "No args throws error",
			ids:      []string{},
			provided: nil,
			expected: &model.Work{
				Title:       "",
				Description: "",
				Duration:    0,
				When:        time.Time{},
				Tags:        []string{},
			},
			err: errors.New(e.EditID),
		}, {
			name:     "2 ids throws error",
			ids:      []string{id, id},
			provided: nil,
			expected: nil,
			err:      errors.New(e.EditID),
		}, {
			name: "Empty parameters default correctly",
			ids:  []string{id},
			provided: &model.Work{
				Title:       "",
				Description: "",
				Duration:    -1,
				When:        time.Time{},
				Tags:        []string{},
			},
			expected: &model.Work{
				Title:       "",
				Description: "",
				Duration:    -1,
				When:        now,
				Tags:        []string{},
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			if testItem.provided != nil {
				var whenString string
				if (testItem.provided.When == time.Time{}) {
					whenString = helpers.TimeFormat(now)
				} else {
					whenString = helpers.TimeFormat(testItem.provided.When)
				}
				setProvidedEditValues(
					testItem.provided.Title,
					testItem.provided.Description,
					testItem.provided.Duration,
					testItem.provided.Author,
					whenString,
					strings.Join(testItem.provided.Tags, ", "))
			}
			retErr := editArgs(testItem.ids)

			if testItem.err != nil {
				assert.EqualError(t, retErr, testItem.err.Error())
			} else {
				assert.Nil(t, retErr)

				if testItem.provided != nil {
					assert.Equal(t, testItem.expected.Title, editTitle)
					assert.Equal(t, testItem.expected.Description, editDescription)
					assert.Equal(t, testItem.expected.Duration, editDuration)
					assert.Equal(t, testItem.expected.Author, editAuthor)
					assert.Equal(t, testItem.expected.When, editWhen)
					assert.Equal(t, testItem.expected.Tags, editTags)
				}
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
