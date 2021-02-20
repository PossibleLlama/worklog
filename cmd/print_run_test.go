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
		svcMethod   string
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
			svcMethod: "GetWorklogsBetween",
			expErr:    nil,
		}, {
			name:        "With only start date",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			author:      helpers.RandString(shortLength),
			tags: []string{
				helpers.RandString(shortLength),
				helpers.RandString(shortLength)},
			ids:       []string{},
			startDate: testDefaultStartDate,
			endDate:   time.Time{},
			svcMethod: "GetWorklogsBetween",
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
			svcMethod: "GetWorklogsByID",
			expErr:    nil,
		}, {
			name:        "With both date and ID",
			title:       helpers.RandString(shortLength),
			description: helpers.RandString(shortLength),
			author:      helpers.RandString(shortLength),
			tags: []string{
				helpers.RandString(shortLength),
				helpers.RandString(shortLength)},
			ids: []string{
				helpers.RandString(shortLength)},
			startDate: testDefaultStartDate,
			endDate:   testDefaultEndDate,
			svcMethod: "GetWorklogsByID",
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
			svcMethod: "GetWorklogsBetween",
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
		if testItem.svcMethod == "GetWorklogsBetween" {
			mockService.
				On(testItem.svcMethod,
					testItem.startDate,
					testItem.endDate,
					expFilter).
				Return(expWl,
					0,
					testItem.expErr)
		} else {
			mockService.
				On(testItem.svcMethod,
					expFilter,
					testItem.ids).
				Return(expWl,
					0,
					testItem.expErr)
		}
		wlService = mockService

		t.Run(testItem.name, func(t *testing.T) {
			actualErr := printRun(testItem.ids...)

			mockService.AssertExpectations(t)

			if testItem.svcMethod == "GetWorklogsBetween" {
				mockService.AssertCalled(t,
					testItem.svcMethod,
					testItem.startDate,
					testItem.endDate,
					expFilter)
			} else {
				mockService.
					On(testItem.svcMethod,
						expFilter,
						testItem.ids).
					Return(expWl,
						0,
						testItem.expErr)
			}
			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
