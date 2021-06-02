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

func setProvidedEditValuesRun(id, title, description string, duration int, author string, when time.Time, tags []string) {
	setProvidedEditValues(title, description, duration, author, "", "")
	editID = id
	editWhen = when
	editTags = tags
}

func TestEditArgs(t *testing.T) {
	id := helpers.RandAlphabeticString(shortLength)
	title := helpers.RandHexAlphaNumericString(shortLength)
	description := helpers.RandHexAlphaNumericString(longLength)
	now := time.Date(2020, time.January, 30, 23, 59, 0, 0, time.UTC)
	tag := helpers.RandHexAlphaNumericString(longLength)

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
		}, {
			name: "Filled parameters are used correctly",
			ids:  []string{id},
			provided: &model.Work{
				Title:       title,
				Description: description,
				Duration:    longLength,
				When:        time.Date(2021, time.February, 10, 21, 56, 45, 0, time.UTC),
				Tags:        []string{tag},
			},
			expected: &model.Work{
				Title:       title,
				Description: description,
				Duration:    longLength,
				When:        time.Date(2021, time.February, 10, 21, 56, 45, 0, time.UTC),
				Tags:        []string{tag},
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
	now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	tm := time.Date(2021, time.February, 10, 21, 56, 45, 0, time.UTC)

	var tests = []struct {
		name     string
		provided model.Work
		err      error
	}{
		{
			name: "Sends model to service",
			provided: model.Work{
				ID:          helpers.RandAlphabeticString(shortLength),
				Revision:    1,
				Title:       helpers.RandAlphabeticString(shortLength),
				Description: helpers.RandAlphabeticString(longLength),
				Author:      helpers.RandAlphabeticString(shortLength),
				Duration:    longLength,
				When:        tm,
				Tags: []string{
					helpers.RandAlphabeticString(shortLength),
					helpers.RandAlphabeticString(shortLength),
				},
				CreatedAt: now,
			},
			err: nil,
		}, {
			name: "Error passed back",
			provided: model.Work{
				ID:          helpers.RandAlphabeticString(shortLength),
				Revision:    1,
				Title:       helpers.RandAlphabeticString(shortLength),
				Description: helpers.RandAlphabeticString(longLength),
				Author:      helpers.RandAlphabeticString(shortLength),
				Duration:    longLength,
				When:        tm,
				Tags: []string{
					helpers.RandAlphabeticString(shortLength),
					helpers.RandAlphabeticString(shortLength),
				},
				CreatedAt: now,
			},
			err: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		mockService := new(service.MockService)
		mockService.On("EditWorklog", testItem.provided.ID, &testItem.provided).Return(0, testItem.err)
		wlService = mockService

		t.Run(testItem.name, func(t *testing.T) {
			setProvidedEditValuesRun(
				testItem.provided.ID,
				testItem.provided.Title,
				testItem.provided.Description,
				testItem.provided.Duration,
				testItem.provided.Author,
				testItem.provided.When,
				testItem.provided.Tags)
			retErr := editRun([]string{})

			mockService.AssertExpectations(t)
			mockService.AssertCalled(t,
				"EditWorklog",
				testItem.provided.ID,
				&testItem.provided)
			assert.Equal(t, testItem.err, retErr)
		})
	}
}
