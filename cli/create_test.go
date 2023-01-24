package cli

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
	randTitle := helpers.RandAlphabeticString(shortLength)
	randDesc := helpers.RandAlphabeticString(shortLength)
	randAuth := helpers.RandAlphabeticString(shortLength)
	randTagA := helpers.RandAlphabeticString(shortLength)
	randTagB := helpers.RandAlphabeticString(shortLength)

	var tests = []struct {
		name             string
		inputTitle       string
		expTitle         string
		inputDescription string
		expDescription   string
		inputAuthor      string
		expAuthor        string
		inputDuration    int
		expDuration      int
		inputTags        string
		expTags          []string
		inputWhen        string
		expWhen          time.Time
		expErr           error
	}{
		{
			name:             "Uses inputs",
			inputTitle:       randTitle,
			expTitle:         randTitle,
			inputDescription: randDesc,
			expDescription:   randDesc,
			inputAuthor:      randAuth,
			expAuthor:        randAuth,
			inputDuration:    shortLength,
			expDuration:      shortLength,
			inputTags:        randTagA + "," + randTagB,
			expTags:          []string{randTagA, randTagB},
			inputWhen:        now.Format(time.RFC3339),
			expWhen:          now,
			expErr:           nil,
		}, {
			name:             "Negative duration",
			inputTitle:       randTitle,
			expTitle:         randTitle,
			inputDescription: randDesc,
			expDescription:   randDesc,
			inputAuthor:      randAuth,
			expAuthor:        randAuth,
			inputDuration:    -1,
			expDuration:      -1,
			inputTags:        randTagA + "," + randTagB,
			expTags:          []string{randTagA, randTagB},
			inputWhen:        now.Format(time.RFC3339),
			expWhen:          now,
			expErr:           nil,
		}, {
			name:             "Empty author stays empty",
			inputTitle:       randTitle,
			expTitle:         randTitle,
			inputDescription: randDesc,
			expDescription:   randDesc,
			inputAuthor:      "",
			expAuthor:        "",
			inputDuration:    shortLength,
			expDuration:      shortLength,
			inputTags:        randTagA + "," + randTagB,
			expTags:          []string{randTagA, randTagB},
			inputWhen:        now.Format(time.RFC3339),
			expWhen:          now,
			expErr:           nil,
		}, {
			name:             "Empty tags",
			inputTitle:       randTitle,
			expTitle:         randTitle,
			inputDescription: randDesc,
			expDescription:   randDesc,
			inputAuthor:      randAuth,
			expAuthor:        randAuth,
			inputDuration:    shortLength,
			expDuration:      shortLength,
			inputTags:        "",
			expTags:          []string{},
			inputWhen:        now.Format(time.RFC3339),
			expWhen:          now,
			expErr:           nil,
		}, {
			name:             "Whitespace tags",
			inputTitle:       randTitle,
			expTitle:         randTitle,
			inputDescription: randDesc,
			expDescription:   randDesc,
			inputAuthor:      randAuth,
			expAuthor:        randAuth,
			inputDuration:    shortLength,
			expDuration:      shortLength,
			inputTags:        " ,\t",
			expTags:          []string{},
			inputWhen:        now.Format(time.RFC3339),
			expWhen:          now,
			expErr:           nil,
		}, {
			name:             "Single character tag",
			inputTitle:       randTitle,
			expTitle:         randTitle,
			inputDescription: randDesc,
			expDescription:   randDesc,
			inputAuthor:      randAuth,
			expAuthor:        randAuth,
			inputDuration:    shortLength,
			expDuration:      shortLength,
			inputTags:        "1",
			expTags:          []string{"1"},
			inputWhen:        now.Format(time.RFC3339),
			expWhen:          now,
			expErr:           nil,
		}, {
			name:             "Invalid when",
			inputTitle:       randTitle,
			expTitle:         randTitle,
			inputDescription: randDesc,
			expDescription:   randDesc,
			inputAuthor:      randAuth,
			expAuthor:        randAuth,
			inputDuration:    shortLength,
			expDuration:      shortLength,
			inputTags:        randTagA + "," + randTagB,
			expTags:          []string{randTagA, randTagB},
			inputWhen:        helpers.RandAlphabeticString(shortLength),
			expWhen:          time.Time{},
			expErr:           errors.New("Could not find format for \""),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedCreateArgValues(
				testItem.inputTitle,
				testItem.inputDescription,
				testItem.inputAuthor,
				testItem.inputWhen,
				testItem.inputDuration,
				testItem.inputTags)

			actualErr := createArgs()

			if testItem.expErr == nil {
				assert.Nil(t, actualErr)
			} else {
				assert.Contains(t, actualErr.Error(), testItem.expErr.Error())
			}

			assert.Equal(t, testItem.expTitle, createTitle)
			assert.Equal(t, testItem.expDescription, createDescription)
			assert.Equal(t, testItem.expAuthor, createAuthor)
			assert.Equal(t, testItem.expDuration, createDuration)
			assert.Equal(t, testItem.inputTags, createTagsString)
			assert.Equal(t, testItem.expTags, createTags)
			assert.Equal(t, testItem.inputWhen, createWhenString)
			assert.Equal(t, testItem.expWhen, createWhen)
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
		assert.Nil(t, actualErr, "An error occurred when creating without args")

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
		assert.Nil(t, actualErr, "An error occurred when creating with args")

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
				ID:             id,
				Revision:       1,
				Title:          testItem.title,
				Description:    testItem.description,
				Author:         defaultAuthor,
				Duration:       testItem.duration,
				Tags:           testItem.tags,
				When:           now,
				WhenQueryEpoch: now.Unix(),
				CreatedAt:      now,
			}
		} else {
			w = &model.Work{
				ID:             id,
				Revision:       1,
				Title:          testItem.title,
				Description:    testItem.description,
				Author:         testItem.author,
				Duration:       testItem.duration,
				Tags:           testItem.tags,
				When:           now,
				WhenQueryEpoch: now.Unix(),
				CreatedAt:      now,
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
