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
var testDefaultFilter = &model.Work{
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

func setProvidedPrintArgValues(title, description, author string, fr format, tags []string, s, e string) {
	setFormatValues(fr)

	printFilterTitle = title
	printFilterDescription = description
	printFilterAuthor = author
	printFilterTagsString = strings.Join(tags, ",")
	printFilterTags = []string{}

	printStartDate = time.Time{}
	printStartDateString = s
	printEndDate = time.Time{}
	printEndDateString = e

	printToday = false
	printThisWeek = false
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
			testDefaultFilter.Title,
			testDefaultFilter.Description,
			testDefaultFilter.Author,
			testItem.usedFormat,
			testDefaultFilter.Tags,
			testDefaultStartDate.Format(time.RFC3339),
			testDefaultEndDate.Format(time.RFC3339))

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
		usedFilter *model.Work
		expFilter  *model.Work
	}{
		{
			name: "No filters",
			usedFilter: &model.Work{
				Title:       "",
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
			expFilter: &model.Work{
				Title:       "",
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
		}, {
			name: "Single filter",
			usedFilter: &model.Work{
				Title:       randString,
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
			expFilter: &model.Work{
				Title:       randString,
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
		}, {
			name: "Multiple filters",
			usedFilter: &model.Work{
				Title:       randString,
				Description: "",
				Author:      "",
				Tags:        []string{randString},
			},
			expFilter: &model.Work{
				Title:       randString,
				Description: "",
				Author:      "",
				Tags:        []string{randString},
			},
		}, {
			name: "Single filter with postfix spacing",
			usedFilter: &model.Work{
				Title: randString + " ",
				Tags:  []string{},
			},
			expFilter: &model.Work{
				Title: randString,
				Tags:  []string{},
			},
		}, {
			name: "Single filter with prefix spacing",
			usedFilter: &model.Work{
				Title: " " + randString,
				Tags:  []string{},
			},
			expFilter: &model.Work{
				Title: randString,
				Tags:  []string{},
			},
		}, {
			name: "Empty string for tag does not filter",
			usedFilter: &model.Work{
				Tags: []string{""},
			},
			expFilter: &model.Work{
				Tags: []string{},
			},
		},
	}

	for _, testItem := range tests {
		setProvidedPrintArgValues(
			testItem.usedFilter.Title,
			testItem.usedFilter.Description,
			testItem.usedFilter.Author,
			testDefaultFormat,
			testItem.usedFilter.Tags,
			testDefaultStartDate.Format(time.RFC3339),
			testDefaultEndDate.Format(time.RFC3339))

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

	var tests = []struct {
		name     string
		sDate    string
		expSDate time.Time
		eDate    string
		expEDate time.Time
		expErr   error
	}{
		{
			name:     "Start date is required before end date is used",
			sDate:    "",
			expSDate: time.Time{},
			eDate:    testDefaultEndDate.Format(time.RFC3339),
			expErr:   errors.New("one flag is required"),
		}, {
			name:     "Invalid string for start date throws error",
			sDate:    randString,
			expSDate: time.Time{},
			eDate:    "",
			expErr:   errors.New("unable to parse string as date. 'parsing time \"" + randString + "\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"" + randString + "\" as \"2006\"'"),
		}, {
			name:     "End date is not required if start date provided",
			sDate:    testDefaultStartDate.Format(time.RFC3339),
			expSDate: testDefaultStartDate,
			eDate:    "",
			expErr:   nil,
		}, {
			name:     "Invalid string for end date throws error when start date provided",
			sDate:    testDefaultStartDate.Format(time.RFC3339),
			expSDate: testDefaultStartDate,
			eDate:    randString,
			expErr:   errors.New("unable to parse string as date. 'parsing time \"" + randString + "\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"" + randString + "\" as \"2006\"'"),
		},
	}

	for _, testItem := range tests {
		setProvidedPrintArgValues(
			testDefaultFilter.Title,
			testDefaultFilter.Description,
			testDefaultFilter.Author,
			testDefaultFormat,
			testDefaultFilter.Tags,
			testItem.sDate,
			testItem.eDate)

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
