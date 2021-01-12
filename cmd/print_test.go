package cmd

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/assert"
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
	providedStartDate = time.Date(y1, m1, d1, 06, 00, 00, 00, time.UTC)
	expectedStartDate = time.Date(y1, m1, d1, 00, 00, 00, 00, time.UTC)
	providedEndDate = time.Date(y2, m2, d2, 12, 00, 00, 00, time.UTC)
	expectedEndDate = time.Date(y3, m3, d3, 00, 00, 00, 00, time.UTC)

	var tests = []struct {
		name       string
		usedFormat format
		expFormat  format
		filter     *model.Work
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
			filter: &model.Work{
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
			filter: &model.Work{
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
			filter: &model.Work{
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
			filter: &model.Work{
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
			filter: &model.Work{
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
			filter: &model.Work{
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
			filter: &model.Work{
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
			filter: &model.Work{
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
			filter: &model.Work{
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
			filter: &model.Work{
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
				yaml:   false,
				json:   false,
			},
			expFormat: format{
				pretty: true,
				yaml:   false,
				json:   false,
			},
			filter: &model.Work{
				Title:       helpers.RandString(shortLength) + " ",
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Single filter with prefix spacing",
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
			filter: &model.Work{
				Title:       " " + helpers.RandString(shortLength),
				Description: "",
				Author:      "",
				Tags:        []string{},
			},
			sDate:  providedStartDate.Format(time.RFC3339),
			eDate:  providedEndDate.Format(time.RFC3339),
			expErr: nil,
		},
	}

	for _, testItem := range tests {
		setProvidedPrintArgValues(
			testItem.filter.Title,
			testItem.filter.Description,
			testItem.filter.Author,
			testItem.usedFormat,
			testItem.filter.Tags,
			testItem.sDate,
			testItem.eDate)
		t.Run(testItem.name, func(t *testing.T) {
			actualErr := printArgs()

			assert.Equal(t, strings.TrimSpace(testItem.filter.Title), printFilterTitle)
			assert.Equal(t, strings.TrimSpace(testItem.filter.Description), printFilterDescription)
			assert.Equal(t, strings.TrimSpace(testItem.filter.Author), printFilterAuthor)
			assert.Equal(t, strings.Join(testItem.filter.Tags, ","), printFilterTagsString)
			assert.Equal(t, testItem.filter.Tags, printFilterTags)

			assert.Equal(t, testItem.expFormat.pretty, printOutputPretty)
			assert.Equal(t, testItem.expFormat.yaml, printOutputYAML)
			assert.Equal(t, testItem.expFormat.json, printOutputJSON)

			assert.Equal(t, expectedStartDate, printStartDate, fmt.Sprintf("Start: Exp: %s, Act: %s", expectedStartDate, printStartDate))
			assert.Equal(t, testItem.sDate, printStartDateString)
			assert.Equal(t, expectedEndDate, printEndDate, fmt.Sprintf("End: Exp: %s, Act: %s", expectedEndDate, printEndDate))
			assert.Equal(t, testItem.eDate, printEndDateString)

			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
