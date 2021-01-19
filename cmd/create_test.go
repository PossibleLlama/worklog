package cmd

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/service"
	"github.com/stretchr/testify/assert"
)

var (
	defaultDuration = 0
	defaultAuthor   = ""
)

func setProvidedCreateArgValues(title, description, when string, duration int, tags string) {
	createTitle = title
	createDescription = description
	createDuration = duration
	createTagsString = tags
	createTags = []string{}
	createWhenString = when
	createWhen = time.Time{}
}

func setProvidedCreateRunValues(id, title, description string, when time.Time, duration int, tags []string) {
	createID = id
	createTitle = title
	createDescription = description
	createDuration = duration
	createTagsString = strings.Join(tags, ",")
	createTags = tags
	createWhenString = when.Format(time.RFC3339)
	createWhen = when
}

func TestCreateArgs(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	var tests = []struct {
		name        string
		title       string
		description string
		duration    int
		tagsString  string
		tags        []string
		whenString  string
		when        time.Time
		expErr      error
	}{
		{
			name:        "Variables take defaults",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			duration:    shortLength,
			tagsString:  "alpha, beta",
			tags:        []string{"alpha", "beta"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Negative duration",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			duration:    -1,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Padded title",
			title:       "\n" + helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			duration:    -1,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Padded description",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength) + " ",
			duration:    longLength,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Padded when",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			duration:    longLength,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  "\t" + now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Invalid when",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			duration:    longLength,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  helpers.RandString(shortLength),
			when:        time.Time{},
			expErr:      errors.New("unable to parse string as date. 'parsing time"),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedCreateArgValues(
				testItem.title,
				testItem.description,
				testItem.whenString,
				testItem.duration,
				testItem.tagsString)

			actualErr := createArgs()

			if testItem.expErr == nil {
				assert.Nil(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), testItem.expErr.Error())
			}

			assert.Equal(t, strings.TrimSpace(testItem.title), createTitle)
			assert.Equal(t, strings.TrimSpace(testItem.description), createDescription)

			if testItem.duration >= 0 {
				assert.Equal(t, testItem.duration, createDuration)
			} else {
				assert.Equal(t, 0, createDuration)
			}
			assert.Equal(t, testItem.tagsString, createTagsString)
			assert.Equal(t, testItem.tags, createTags)
			assert.Equal(t, testItem.whenString, createWhenString)
			assert.Equal(t, testItem.when, createWhen)
		})
	}
}

func TestCreateRun(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	var tests = []struct {
		name        string
		title       string
		description string
		duration    int
		tags        []string
		when        time.Time
		expErr      error
	}{
		{
			name:        "Send to service",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			duration:    longLength,
			tags: []string{
				helpers.RandString(shortLength),
				helpers.RandString(shortLength)},
			when:   now,
			expErr: nil,
		}, {
			name:        "Error passed back",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			duration:    longLength,
			tags: []string{
				helpers.RandString(shortLength),
				helpers.RandString(shortLength)},
			when:   now,
			expErr: nil,
		},
	}

	for _, testItem := range tests {
		id := helpers.RandString(shortLength)
		w := &model.Work{
			ID:          id,
			Revision:    1,
			Title:       testItem.title,
			Description: testItem.description,
			Author:      defaultAuthor,
			Duration:    testItem.duration,
			Where:       "",
			Tags:        testItem.tags,
			When:        now,
			CreatedAt:   now,
		}
		mockService := new(service.MockService)
		mockService.On("CreateWorklog", w).Return(0, testItem.expErr)
		wlService = mockService

		t.Run(testItem.name, func(t *testing.T) {
			setProvidedCreateRunValues(
				id,
				testItem.title,
				testItem.description,
				testItem.when,
				testItem.duration,
				testItem.tags)

			actualErr := createRun()

			mockService.AssertExpectations(t)
			mockService.AssertCalled(t, "CreateWorklog", w)
			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
