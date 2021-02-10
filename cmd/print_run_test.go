package cmd

import (
	"errors"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/service"
	"github.com/stretchr/testify/assert"
)

func setProvidedPrintRunValues(fr format, title, description, author string, tags []string, startDate, endDate time.Time) {
	setFormatValues(fr)

	printFilterTitle = title
	printFilterDescription = description
	printFilterAuthor = author
	printFilterTagsString = ""
	printFilterTags = tags

	printStartDate = startDate
	printStartDateString = ""
	printEndDate = endDate
	printEndDateString = ""

	printToday = false
	printThisWeek = false
}

func TestPrintRun(t *testing.T) {
	expWl := []*model.Work{model.NewWork(
		helpers.RandString(shortLength),
		helpers.RandString(shortLength),
		helpers.RandString(shortLength),
		shortLength,
		[]string{helpers.RandString(shortLength)},
		time.Time{})}

	var tests = []struct {
		name        string
		title       string
		description string
		author      string
		tags        []string
		ids         []string
		startDate   time.Time
		endDate     time.Time
		expErr      error
	}{
		{
			name:        "With only date",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			author:      helpers.RandString(shortLength),
			tags: []string{
				helpers.RandString(shortLength),
				helpers.RandString(shortLength)},
			ids:       []string{},
			startDate: testDefaultStartDate,
			endDate:   testDefaultEndDate,
			expErr:    nil,
		}, {
			name:        "With only ID",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			author:      helpers.RandString(shortLength),
			tags: []string{
				helpers.RandString(shortLength),
				helpers.RandString(shortLength)},
			ids: []string{
				helpers.RandString(shortLength)},
			startDate: time.Time{},
			endDate:   time.Time{},
			expErr:    nil,
		}, {
			name:        "Error",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			author:      helpers.RandString(shortLength),
			tags: []string{
				helpers.RandString(shortLength),
				helpers.RandString(shortLength)},
			ids:       []string{},
			startDate: testDefaultStartDate,
			endDate:   testDefaultEndDate,
			expErr:    errors.New(helpers.RandString(shortLength)),
		},
	}

	for _, testItem := range tests {
		setProvidedPrintRunValues(format{
			pretty: true,
			yaml:   false,
			json:   false,
		},
			testItem.title,
			testItem.description,
			testItem.author,
			testItem.tags,
			testItem.startDate,
			testItem.endDate)

		expFilter := &model.Work{
			Title:       testItem.title,
			Description: testItem.description,
			Author:      testItem.author,
			Tags:        testItem.tags,
			Duration:    -1,
			When:        time.Time{},
			CreatedAt:   time.Time{},
		}

		mockService := new(service.MockService)
		mockService.
			On("GetWorklogsBetween",
				testItem.startDate,
				testItem.endDate,
				expFilter).
			Return(expWl,
				0,
				testItem.expErr)
		wlService = mockService

		t.Run(testItem.name, func(t *testing.T) {
			actualErr := printRun(testItem.ids...)

			mockService.AssertExpectations(t)
			mockService.AssertCalled(t,
				"GetWorklogsBetween",
				testItem.startDate,
				testItem.endDate,
				expFilter)
			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
