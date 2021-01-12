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

const (
	zero = 00
)

var (
	providedStartDate time.Time
	expectedStartDate time.Time
	providedEndDate   time.Time
	expectedEndDate   time.Time
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

func TestPrintArgs(t *testing.T) {
	now := time.Now()
	y1, m1, d1 := now.Date()
	y2, m2, d2 := now.Add(time.Hour * 24).Date()
	y3, m3, d3 := now.Add(time.Hour * 48).Date()
	providedStartDate = time.Date(y1, m1, d1, 06, zero, zero, zero, time.UTC)
	expectedStartDate = time.Date(y1, m1, d1, zero, zero, zero, zero, time.UTC)
	providedEndDate = time.Date(y2, m2, d2, 12, zero, zero, zero, time.UTC)
	expectedEndDate = time.Date(y3, m3, d3, zero, zero, zero, zero, time.UTC)

	randString := helpers.RandString(shortLength)

	var tests = []struct {
		name       string
		usedFormat format
		expFormat  format
		usedFilter *model.Work
		expFilter  *model.Work
		sDate      string
		eDate      string
		expErr     error
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: helpers.RandString(shortLength),
				Author:      helpers.RandString(shortLength),
				Tags: []string{
					helpers.RandString(shortLength),
					helpers.RandString(shortLength)},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: helpers.RandString(shortLength),
				Author:      helpers.RandString(shortLength),
				Tags: []string{
					helpers.RandString(shortLength),
					helpers.RandString(shortLength)},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: helpers.RandString(shortLength),
				Author:      helpers.RandString(shortLength),
				Tags: []string{
					helpers.RandString(shortLength),
					helpers.RandString(shortLength)},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: helpers.RandString(shortLength),
				Author:      helpers.RandString(shortLength),
				Tags: []string{
					helpers.RandString(shortLength),
					helpers.RandString(shortLength)},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: helpers.RandString(shortLength),
				Author:      helpers.RandString(shortLength),
				Tags: []string{
					helpers.RandString(shortLength),
					helpers.RandString(shortLength)},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: helpers.RandString(shortLength),
				Author:      helpers.RandString(shortLength),
				Tags: []string{
					helpers.RandString(shortLength),
					helpers.RandString(shortLength)},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: helpers.RandString(shortLength),
				Author:      helpers.RandString(shortLength),
				Tags: []string{
					helpers.RandString(shortLength),
					helpers.RandString(shortLength)},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "No filters",
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
			usedFilter: &model.Work{
				Title:       "",
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Single filter",
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Multiple filters",
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
			usedFilter: &model.Work{
				Title:       helpers.RandString(shortLength),
				Description: "",
				Author:      "",
				Tags:        []string{helpers.RandString(shortLength)},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Single filter with postfix spacing",
			usedFormat: format{
				pretty: true,
			},
			expFormat: format{
				pretty: true,
			},
			usedFilter: &model.Work{
				Title: randString + " ",
				Tags:  []string{},
			},
			expFilter: &model.Work{
				Title: randString,
				Tags:  []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Single filter with prefix spacing",
			usedFormat: format{
				pretty: true,
			},
			expFormat: format{
				pretty: true,
			},
			usedFilter: &model.Work{
				Title: " " + randString,
				Tags:  []string{},
			},
			expFilter: &model.Work{
				Title: randString,
				Tags:  []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Empty string for tag does not filter",
			usedFormat: format{
				pretty: true,
			},
			expFormat: format{
				pretty: true,
			},
			usedFilter: &model.Work{
				Tags: []string{""},
			},
			expFilter: &model.Work{
				Tags: []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Empty string for start date throws error",
			usedFormat: format{
				pretty: true,
			},
			expFormat: format{
				pretty: true,
			},
			usedFilter: &model.Work{
				Tags: []string{},
			},
			expFilter: &model.Work{
				Tags: []string{},
			},
			sDate:  "",
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: errors.New("one flag is required"),
		}, {
			name: "Invalid string for start date throws error",
			usedFormat: format{
				pretty: true,
			},
			expFormat: format{
				pretty: true,
			},
			usedFilter: &model.Work{
				Tags: []string{},
			},
			expFilter: &model.Work{
				Tags: []string{},
			},
			sDate:  randString,
			eDate:  "",
			expErr: errors.New("unable to parse string as date. 'parsing time \"" + randString + "\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"" + randString + "\" as \"2006\"'"),
		}, {
			name: "Empty string for end date with valid date does not throw error",
			usedFormat: format{
				pretty: true,
			},
			expFormat: format{
				pretty: true,
			},
			usedFilter: &model.Work{
				Tags: []string{},
			},
			expFilter: &model.Work{
				Tags: []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  "",
			expErr: nil,
		}, {
			name: "Invalid string for end date throws error",
			usedFormat: format{
				pretty: true,
			},
			expFormat: format{
				pretty: true,
			},
			usedFilter: &model.Work{
				Tags: []string{},
			},
			expFilter: &model.Work{
				Tags: []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  randString,
			expErr: errors.New("unable to parse string as date. 'parsing time \"" + randString + "\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"" + randString + "\" as \"2006\"'"),
		},
	}

	for _, testItem := range tests {
		setProvidedPrintArgValues(
			testItem.usedFilter.Title,
			testItem.usedFilter.Description,
			testItem.usedFilter.Author,
			testItem.usedFormat,
			testItem.usedFilter.Tags,
			testItem.sDate,
			testItem.eDate)
		t.Run(testItem.name, func(t *testing.T) {
			actualErr := printArgs()

			if testItem.expFilter == nil {
				assert.Equal(t, testItem.usedFilter.Title, printFilterTitle)
				assert.Equal(t, testItem.usedFilter.Description, printFilterDescription)
				assert.Equal(t, testItem.usedFilter.Author, printFilterAuthor)
				assert.Equal(t, strings.Join(testItem.usedFilter.Tags, ","), printFilterTagsString)
				assert.Equal(t, testItem.usedFilter.Tags, printFilterTags)
			} else {
				assert.Equal(t, testItem.expFilter.Title, printFilterTitle)
				assert.Equal(t, testItem.expFilter.Description, printFilterDescription)
				assert.Equal(t, testItem.expFilter.Author, printFilterAuthor)
				assert.Equal(t, strings.Join(testItem.usedFilter.Tags, ","), printFilterTagsString)
				assert.Equal(t, testItem.expFilter.Tags, printFilterTags)
			}

			assert.Equal(t, testItem.expFormat.pretty, printOutputPretty)
			assert.Equal(t, testItem.expFormat.yaml, printOutputYAML)
			assert.Equal(t, testItem.expFormat.json, printOutputJSON)

			if len(testItem.sDate) == 20 {
				assert.Equal(t, expectedStartDate, printStartDate, fmt.Sprintf("Start: Exp: %s, Act: %s", expectedStartDate, printStartDate))
				assert.Equal(t, testItem.sDate, printStartDateString)
				if len(testItem.eDate) == 20 {
					assert.Equal(t, expectedEndDate, printEndDate, fmt.Sprintf("End: Exp: %s, Act: %s", expectedEndDate, printEndDate))
					assert.Equal(t, testItem.eDate, printEndDateString)
				}
			}
			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
