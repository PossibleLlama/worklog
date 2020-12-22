package model

import (
	"testing"
	"time"
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
		expected     Work
	}{
		{
			"Full work",
			"title",
			"description",
			"who",
			15,
			[]string{"alpha", "beta"},
			validDate,
			Work{
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
			expected: Work{
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
			before := time.Now()
			actual := NewWork(
				testItem.wTitle,
				testItem.wDescription,
				testItem.wAuthor,
				testItem.wDuration,
				testItem.wTags,
				testItem.wWhen)
			after := time.Now()

			if actual.Title != testItem.expected.Title {
				t.Errorf("Should have title %s, instead has %s", testItem.expected.Title, actual.Title)
			}
			if actual.Description != testItem.expected.Description {
				t.Errorf("Should have description %s, instead has %s", testItem.expected.Description, actual.Description)
			}
			if actual.Author != testItem.expected.Author {
				t.Errorf("Should have author %s, instead has %s", testItem.expected.Author, actual.Author)
			}
			if actual.Where != testItem.expected.Where {
				t.Errorf("Should have where %s, instead has %s", testItem.expected.Where, actual.Where)
			}
			if actual.Duration != testItem.expected.Duration {
				t.Errorf("Should have duration %d, instead has %d", testItem.expected.Duration, actual.Duration)
			}
			if len(actual.Tags) != len(testItem.expected.Tags) {
				t.Errorf("Should have same number of tags %d, instead has %d", len(testItem.expected.Tags), len(actual.Tags))
			}
			for index, element := range testItem.expected.Tags {
				if actual.Tags[index] != element {
					t.Errorf("Item %s does not match element %s in list %s",
						element, actual.Tags[index], actual.Tags)
				}
			}
			if !actual.When.Equal(testItem.expected.When) {
				t.Errorf("Should have when '%s', instead has '%s'", actual.When, actual.Created)
			}
			if time.Since(actual.Created) < time.Since(before) {
				t.Error("Was not created after start of test")
			}
			if time.Since(actual.Created) < time.Since(after) {
				t.Error("Was not created before end of test")
			}
		})
	}
}
