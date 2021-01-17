package cmd

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/assert"
)

const zero = 00

var testDefaultFormat = format{
	pretty: true,
	yaml:   false,
	json:   false,
}
var testDefaultFilter = model.Work{
	Title:       "",
	Description: "",
	Author:      "",
	Tags:        []string{},
}

var (
	now                  = time.Now()
	day                  = time.Hour * 24
	testDefaultStartDate = time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		zero, zero, zero, zero, time.UTC)
	testDefaultEndDate = time.Date(
		now.Add(day).Year(),
		now.Add(day).Month(),
		now.Add(day).Day(),
		zero, zero, zero, zero, time.UTC)
)

type format struct {
	pretty bool
	yaml   bool
	json   bool
}

func setProvidedPrintArgValues(w model.Work, fr format, s, e string, today, week bool) {
	setFormatValues(fr)

	printFilterTitle = w.Title
	printFilterDescription = w.Description
	printFilterAuthor = w.Author
	printFilterTagsString = strings.Join(w.Tags, ",")
	printFilterTags = []string{}

	printStartDate = time.Time{}
	printStartDateString = s
	printEndDate = time.Time{}
	printEndDateString = e

	printToday = today
	printThisWeek = week
}

func setFormatValues(fr format) {
	printOutputPretty = fr.pretty
	printOutputYAML = fr.yaml
	printOutputJSON = fr.json
}

func TestPrintArgsFormat(t *testing.T) {
	var tests = []struct {
		name       string
		usedFormat format
		expFormat  format
	}{
		{
			name: "Full arguments pretty",
			usedFormat: format{
				pretty: true,
				yaml:   false,
				json:   false,
			},
			expFormat: format{
				pretty: true,
				yaml:   false,
				json:   false,
			},
		}, {
			name: "Full arguments yaml",
			usedFormat: format{
				pretty: false,
				yaml:   true,
				json:   false,
			},
			expFormat: format{
				pretty: false,
				yaml:   true,
				json:   false,
			},
		}, {
			name: "Full arguments json",
			usedFormat: format{
				pretty: false,
				yaml:   false,
				json:   true,
			},
			expFormat: format{
				pretty: false,
				yaml:   false,
				json:   true,
			},
		}, {
			name: "All formats",
			usedFormat: format{
				pretty: true,
				yaml:   true,
				json:   true,
			},
			expFormat: format{
				pretty: true,
				yaml:   false,
				json:   false,
			},
		}, {
			name: "Pretty and yaml formats",
			usedFormat: format{
				pretty: true,
				yaml:   true,
				json:   false,
			},
			expFormat: format{
				pretty: true,
				yaml:   false,
				json:   false,
			},
		}, {
			name: "Pretty and json formats",
			usedFormat: format{
				pretty: true,
				yaml:   false,
				json:   true,
			},
			expFormat: format{
				pretty: true,
				yaml:   false,
				json:   false,
			},
		}, {
			name: "Yaml and json formats",
			usedFormat: format{
				pretty: false,
				yaml:   true,
				json:   true,
			},
			expFormat: format{
				pretty: false,
				yaml:   true,
				json:   false,
			},
		},
	}

	for _, testItem := range tests {
		setProvidedPrintArgValues(
			testDefaultFilter,
			testItem.usedFormat,
			testDefaultStartDate.Format(time.RFC3339),
			testDefaultEndDate.Format(time.RFC3339),
			false,
			false)

		t.Run(testItem.name, func(t *testing.T) {
			err := printArgs()

			assert.Nil(t, err)
			assert.Equal(t, testItem.expFormat.pretty, printOutputPretty)
			assert.Equal(t, testItem.expFormat.yaml, printOutputYAML)
			assert.Equal(t, testItem.expFormat.json, printOutputJSON)
		})
	}
}

func TestPrintArgsFilter(t *testing.T) {
	randString := helpers.RandString(shortLength)

	var tests = []struct {
		name       string
		usedFilter model.Work
		expFilter  model.Work
	}{
		{
			name: "No filters",
			usedFilter: model.Work{
				Title:       "",
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
			expFilter: model.Work{
				Title:       "",
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
		}, {
			name: "Single filter",
			usedFilter: model.Work{
				Title:       randString,
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
			expFilter: model.Work{
				Title:       randString,
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
		}, {
			name: "Multiple filters",
			usedFilter: model.Work{
				Title:       randString,
				Description: "",
				Author:      "",
				Tags:        []string{randString},
			},
			expFilter: model.Work{
				Title:       randString,
				Description: "",
				Author:      "",
				Tags:        []string{randString},
			},
		}, {
			name: "Single filter with postfix spacing",
			usedFilter: model.Work{
				Title: randString + " ",
				Tags:  []string{},
			},
			expFilter: model.Work{
				Title: randString,
				Tags:  []string{},
			},
		}, {
			name: "Single filter with prefix spacing",
			usedFilter: model.Work{
				Title: " " + randString,
				Tags:  []string{},
			},
			expFilter: model.Work{
				Title: randString,
				Tags:  []string{},
			},
		}, {
			name: "Empty string for tag does not filter",
			usedFilter: model.Work{
				Tags: []string{""},
			},
			expFilter: model.Work{
				Tags: []string{},
			},
		},
	}

	for _, testItem := range tests {
		setProvidedPrintArgValues(
			testItem.usedFilter,
			testDefaultFormat,
			testDefaultStartDate.Format(time.RFC3339),
			testDefaultEndDate.Format(time.RFC3339),
			false,
			false)

		t.Run(testItem.name, func(t *testing.T) {
			err := printArgs()

			assert.Nil(t, err)
			assert.Equal(t, testItem.expFilter.Title, printFilterTitle)
			assert.Equal(t, testItem.expFilter.Description, printFilterDescription)
			assert.Equal(t, testItem.expFilter.Author, printFilterAuthor)
			assert.Equal(t, strings.Join(testItem.usedFilter.Tags, ","), printFilterTagsString)
			assert.Equal(t, testItem.expFilter.Tags, printFilterTags)
		})
	}
}

func TestPrintArgsDates(t *testing.T) {
	randString := helpers.RandString(shortLength)

	var testMondayStartDate = helpers.GetPreviousMonday(now)

	var tests = []struct {
		name     string
		sDate    string
		expSDate time.Time
		eDate    string
		expEDate time.Time
		today    bool
		week     bool
		expErr   error
	}{
		{
			name:     "Start date is required before end date is used",
			sDate:    "",
			expSDate: time.Time{},
			eDate:    testDefaultEndDate.Format(time.RFC3339),
			today:    false,
			week:     false,
			expErr:   errors.New("one flag is required"),
		}, {
			name:     "Invalid string for start date throws error",
			sDate:    randString,
			expSDate: time.Time{},
			eDate:    "",
			today:    false,
			week:     false,
			expErr:   errors.New("unable to parse string as date. 'parsing time \"" + randString + "\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"" + randString + "\" as \"2006\"'"),
		}, {
			name:     "End date is not required if start date provided",
			sDate:    testDefaultStartDate.Format(time.RFC3339),
			expSDate: testDefaultStartDate,
			eDate:    "",
			today:    false,
			week:     false,
			expErr:   nil,
		}, {
			name:     "Invalid string for end date throws error when start date provided",
			sDate:    testDefaultStartDate.Format(time.RFC3339),
			expSDate: testDefaultStartDate,
			eDate:    randString,
			today:    false,
			week:     false,
			expErr:   errors.New("unable to parse string as date. 'parsing time \"" + randString + "\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"" + randString + "\" as \"2006\"'"),
		}, {
			name:     "Today sets start and end date",
			sDate:    "",
			expSDate: testDefaultStartDate,
			eDate:    "",
			expEDate: testDefaultEndDate,
			today:    true,
			week:     false,
			expErr:   nil,
		}, {
			name:     "Week sets start and end date",
			sDate:    "",
			expSDate: testMondayStartDate,
			eDate:    "",
			expEDate: testMondayStartDate.Add(7 * day),
			today:    false,
			week:     true,
			expErr:   nil,
		},
	}

	for _, testItem := range tests {
		setProvidedPrintArgValues(
			testDefaultFilter,
			testDefaultFormat,
			testItem.sDate,
			testItem.eDate,
			testItem.today,
			testItem.week)

		t.Run(testItem.name, func(t *testing.T) {
			actualErr := printArgs()

			debugStringStartDate := fmt.Sprintf(
				"Start: Exp: %s, Act: %s",
				testItem.expSDate,
				printStartDate)
			debugStringEndDate := fmt.Sprintf(
				"End: Exp: %s, Act: %s",
				testItem.expEDate,
				printEndDate)

			assert.Equal(t, testItem.expErr, actualErr)

			// Check string values
			assert.Equal(t, testItem.sDate, printStartDateString)
			assert.Equal(t, testItem.eDate, printEndDateString)

			// Check time values
			assert.Equal(t, testItem.expSDate, printStartDate, debugStringStartDate)
			assert.Equal(t, testItem.expEDate, printEndDate, debugStringEndDate)
		})
	}
}
