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
	startDate time.Time
	endDate   time.Time
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
	y, m, d := time.Now().Date()
	startDate, _ = time.Parse(time.RFC3339, fmt.Sprintf("%d-%s-%dT06:00:00Z", y, m, d))
	endDate, _ = time.Parse(time.RFC3339, fmt.Sprintf("%d-%s-%dT12:00:00Z", y, m, d))

	var tests = []struct {
		name   string
		format format
		filter *model.Work
		sDate  string
		eDate  string
		expErr error
	}{
		{
			name: "Full arguments pretty",
			format: format{
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
			sDate:  startDate.Format(time.RFC3339),
			eDate:  endDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Full arguments yaml",
			format: format{
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
			sDate:  startDate.Format(time.RFC3339),
			eDate:  endDate.Format(time.RFC3339),
			expErr: nil,
		}, {
			name: "Full arguments json",
			format: format{
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
			sDate:  startDate.Format(time.RFC3339),
			eDate:  endDate.Format(time.RFC3339),
			expErr: nil,
		},
	}

	for _, testItem := range tests {
		setProvidedPrintArgValues(
			testItem.filter.Title,
			testItem.filter.Description,
			testItem.filter.Author,
			testItem.format,
			testItem.filter.Tags,
			testItem.sDate,
			testItem.eDate)
		t.Run(testItem.name, func(t *testing.T) {
			actualErr := printArgs()

			assert.Equal(t, testItem.filter.Title, printFilterTitle)
			assert.Equal(t, testItem.filter.Description, printFilterDescription)
			assert.Equal(t, testItem.filter.Author, printFilterAuthor)
			assert.Equal(t, strings.Join(testItem.filter.Tags, ","), printFilterTagsString)
			assert.Equal(t, testItem.filter.Tags, printFilterTags)

			assert.Equal(t, testItem.format.pretty, printOutputPretty)
			assert.Equal(t, testItem.format.yaml, printOutputYAML)
			assert.Equal(t, testItem.format.json, printOutputJSON)

			assert.Equal(t, startDate, printStartDate, fmt.Sprintf("Exp: %s, Act: %s", startDate, testItem.sDate))
			assert.Equal(t, testItem.sDate, printStartDateString)
			// assert.Equal(t, endDate, printEndDate)
			assert.Equal(t, testItem.eDate, printEndDateString)

			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
