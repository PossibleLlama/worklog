package main

import (
	"fmt"
	"testing"
	"time"
)

func TestNewWork(t *testing.T) {
	validStringDate := "1970-12-25"
	validDate, err := time.Parse(time.RFC3339, "1970-12-25T00:00:00Z")
	if err != nil {
		t.Error("Unable to parse initial date")
	}

	var tests = []struct{
		name string
		wTitle string
		wDescription string
		wAuthor string
		wWhere string
		wWhen string
		expected Work
	}{
		{
			"Short date",
			"title",
			"description",
			"who",
			"location",
			validStringDate,
			Work{
				Title: "title",
				Description: "description",
				Author: "who",
				Where: "location",
				When: validDate,
			},
		}, {
			"Full date",
			"title",
			"description",
			"who",
			"location",
			fmt.Sprintf("%sT00:00:00Z", validStringDate),
			Work{
				Title: "title",
				Description: "description",
				Author: "who",
				Where: "location",
				When: validDate,
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			before := time.Now()
			actual := New(testItem.wTitle, testItem.wDescription, testItem.wAuthor, testItem.wWhere, testItem.wWhen)
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
			if actual.When != testItem.expected.When {
				t.Errorf("Should have when %s, instead has %s", testItem.expected.When, actual.When)
			}
			if time.Since(before) < time.Since(actual.Created) {
				t.Error("Was not created after start of test")
			}
			if time.Since(actual.Created) < time.Since(after) {
				t.Error("Was not created before end of test")
			}
		})
	}
}