package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
