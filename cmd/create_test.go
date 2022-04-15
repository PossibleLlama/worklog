package cmd

import (
	"errors"
	"fmt"
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

func setProvidedCreateArgValues(title, description, author, when string, duration int, tags string) {
	createTitle = title
	createDescription = description
	createAuthor = author
	createDuration = duration
	createTagsString = tags
	createTags = []string{}
	createWhenString = when
	createWhen = time.Time{}
}

func setProvidedCreateRunValues(id, title, description, author string, when time.Time, duration int, tags []string) {
	createID = id
	createTitle = title
	createDescription = description
	createAuthor = author
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
		author      string
		duration    int
		tagsString  string
		tags        []string
		whenString  string
		when        time.Time
		expErr      error
	}{
		{
			name:        "Variables take defaults",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    shortLength,
			tagsString:  "alpha, beta",
			tags:        []string{"alpha", "beta"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Negative duration",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    -1,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Padded title",
			title:       "\n" + helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    -1,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Padded description",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength) + " ",
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    longLength,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Empty author uses default",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      "",
			duration:    longLength,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Padded author",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength) + " ",
			duration:    longLength,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Empty tags",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    longLength,
			tagsString:  "",
			tags:        []string{},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Whitespace tags",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    longLength,
			tagsString:  " ,\t",
			tags:        []string{},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Single character tags",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    longLength,
			tagsString:  "1, ",
			tags:        []string{"1"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Padded when",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    longLength,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  "\t" + now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "Invalid when",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    longLength,
			tagsString:  "1, 2",
			tags:        []string{"1", "2"},
			whenString:  helpers.RandAlphabeticString(shortLength),
			when:        time.Time{},
			expErr:      errors.New("Could not find format for \""),
		}, {
			name:        "XSS title",
			title:       xssHtmlOpen + helpers.RandAlphabeticString(shortLength) + xssHtmlClose,
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    shortLength,
			tagsString:  "alpha, beta",
			tags:        []string{"alpha", "beta"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "XSS description",
			title:       helpers.RandAlphabeticString(shortLength),
			description: xssHtmlOpen + helpers.RandAlphabeticString(shortLength) + xssHtmlClose,
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    shortLength,
			tagsString:  "alpha, beta",
			tags:        []string{"alpha", "beta"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "XSS author",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      xssHtmlOpen + helpers.RandAlphabeticString(shortLength) + xssHtmlClose,
			duration:    shortLength,
			tagsString:  "alpha, beta",
			tags:        []string{"alpha", "beta"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		}, {
			name:        "XSS tags",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    shortLength,
			tagsString:  xssHtmlOpen + "alpha, beta" + xssHtmlClose,
			tags:        []string{"alpha", "beta"},
			whenString:  now.Format(time.RFC3339),
			when:        now,
			expErr:      nil,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedCreateArgValues(
				testItem.title,
				testItem.description,
				testItem.author,
				testItem.whenString,
				testItem.duration,
				testItem.tagsString)

			actualErr := createArgs()

			if testItem.expErr == nil {
				assert.Nil(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), testItem.expErr.Error())
			}

			assert.Equal(t, helpers.Sanitize(strings.TrimSpace(testItem.title)), createTitle)
			assert.Equal(t, helpers.Sanitize(strings.TrimSpace(testItem.description)), createDescription)
			assert.Equal(t, helpers.Sanitize(strings.TrimSpace(testItem.author)), createAuthor)

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

func TestCreateArgsWhen(t *testing.T) {
	t.Run("Ensure '--when' is the same when using the flag and not", func(t *testing.T) {
		tm := time.Now()

		// This is how the date is initialised in the create function
		setProvidedCreateArgValues(
			"When without specifying",
			"",
			defaultAuthor,
			helpers.TimeFormat(tm),
			defaultDuration,
			"")
		actualErr := createArgs()
		assert.Nil(t, actualErr, "An error occured when creating without args")

		noWhen := createWhen

		zone, _ := tm.Zone()
		setProvidedCreateArgValues(
			"When with specifying",
			"",
			defaultAuthor,
			fmt.Sprintf("%d-%d-%dT%d:%d:%d %s",
				tm.Year(),
				int(tm.Month()),
				tm.Day(),
				tm.Hour(),
				tm.Minute(),
				tm.Second(),
				zone),
			defaultDuration,
			"")
		actualErr = createArgs()
		assert.Nil(t, actualErr, "An error occured when creating with args")

		// TODO I hate timezones
		assert.Equal(t, noWhen.UTC(), createWhen.UTC())
	})
}

func TestCreateRun(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	var tests = []struct {
		name        string
		title       string
		description string
		author      string
		duration    int
		tags        []string
		when        time.Time
		expErr      error
	}{
		{
			name:        "Send to service",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    longLength,
			tags: []string{
				helpers.RandAlphabeticString(shortLength),
				helpers.RandAlphabeticString(shortLength)},
			when:   now,
			expErr: nil,
		}, {
			name:        "Default author is used",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      "",
			duration:    longLength,
			tags: []string{
				helpers.RandAlphabeticString(shortLength),
				helpers.RandAlphabeticString(shortLength)},
			when:   now,
			expErr: nil,
		}, {
			name:        "Error passed back",
			title:       helpers.RandAlphabeticString(shortLength),
			description: helpers.RandAlphabeticString(shortLength),
			author:      helpers.RandAlphabeticString(shortLength),
			duration:    longLength,
			tags: []string{
				helpers.RandAlphabeticString(shortLength),
				helpers.RandAlphabeticString(shortLength)},
			when:   now,
			expErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		id := helpers.RandAlphabeticString(shortLength)
		var w *model.Work
		if testItem.author == "" {
			w = &model.Work{
				ID:          id,
				Revision:    1,
				Title:       testItem.title,
				Description: testItem.description,
				Author:      defaultAuthor,
				Duration:    testItem.duration,
				Tags:        testItem.tags,
				When:        now,
				CreatedAt:   now,
			}
		} else {
			w = &model.Work{
				ID:          id,
				Revision:    1,
				Title:       testItem.title,
				Description: testItem.description,
				Author:      testItem.author,
				Duration:    testItem.duration,
				Tags:        testItem.tags,
				When:        now,
				CreatedAt:   now,
			}
		}
		mockService := new(service.MockService)
		mockService.On("CreateWorklog", w).Return(0, testItem.expErr)
		wlService = mockService

		t.Run(testItem.name, func(t *testing.T) {
			setProvidedCreateRunValues(
				id,
				testItem.title,
				testItem.description,
				testItem.author,
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
