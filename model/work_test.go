package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/stretchr/testify/assert"
)

const (
	shortLength = 30
	longLength  = 256
	dateString  = "2000-01-30T0000:00:00Z"
)

func TestNewWork(t *testing.T) {
	validDate, err := time.Parse(time.RFC3339, "1970-12-25T00:00:00Z")
	if err != nil {
		t.Error("Unable to parse initial date")
	}

	var tests = []struct {
		name         string
		wTitle       string
		wDescription string
		wAuthor      string
		wDuration    int
		wTags        []string
		wWhen        time.Time
		expected     *Work
	}{
		{
			name:         "Full work",
			wTitle:       "title",
			wDescription: "description",
			wAuthor:      "who",
			wDuration:    15,
			wTags:        []string{"alpha", "beta"},
			wWhen:        validDate,
			expected: &Work{
				Title:       "title",
				Description: "description",
				Author:      "who",
				Where:       "",
				Duration:    15,
				Tags:        []string{"alpha", "beta"},
				When:        validDate,
			},
		},
		{
			name:         "Unordered tags become ordered",
			wTitle:       "title",
			wDescription: "description",
			wAuthor:      "who",
			wDuration:    15,
			wTags:        []string{"4", "2", "1", "3"},
			wWhen:        validDate,
			expected: &Work{
				Title:       "title",
				Description: "description",
				Author:      "who",
				Where:       "",
				Duration:    15,
				Tags:        []string{"1", "2", "3", "4"},
				When:        validDate,
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := NewWork(
				testItem.wTitle,
				testItem.wDescription,
				testItem.wAuthor,
				testItem.wDuration,
				testItem.wTags,
				testItem.wWhen)
			finished := time.Now()

			// Instead of mocking time.Now(), just set the result of it to the expected value
			testItem.expected.CreatedAt = actual.CreatedAt

			assert.Equal(t, testItem.expected, actual)
			assert.True(t, finished.Add(time.Second*-1).Before(actual.CreatedAt))
		})
	}
}

func TestString(t *testing.T) {
	date, _ := helpers.GetStringAsDateTime(dateString)
	var tests = []struct {
		name string
		work *Work
		exp  string
	}{
		{
			name: "Full work",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], When: %s, CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing title",
			work: &Work{
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Description: %s, Author: %s, Duration: %d, Tags: [%s], When: %s, CreatedAt: %s",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing description",
			work: &Work{
				Title:    "title",
				Author:   "author",
				Duration: 15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Author: %s, Duration: %d, Tags: [%s], When: %s, CreatedAt: %s",
				"title",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing author",
			work: &Work{
				Title:       "title",
				Description: "description",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Duration: %d, Tags: [%s], When: %s, CreatedAt: %s",
				"title",
				"description",
				15,
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing duration",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Tags: [%s], When: %s, CreatedAt: %s",
				"title",
				"description",
				"author",
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing tags",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				When:        date,
				CreatedAt:   date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, When: %s, CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Empty tags",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags:        []string{},
				When:        date,
				CreatedAt:   date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, When: %s, CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing when",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Basic when",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      time.Time{},
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Missing createdAt",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], When: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Basic createdAt",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: time.Time{},
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], When: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "No fields",
			work: &Work{},
			exp:  "",
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := testItem.work.String()
			assert.Equal(t, testItem.exp, actual)
		})
	}
}
