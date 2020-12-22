package model

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

// Work data model.
// stores information as to what
// work was done by who and when.
type Work struct {
	Title       string
	Description string
	Author      string
	Where       string
	Duration    int
	Tags        []string
	When        time.Time
	Created     time.Time
}

// NewWork is the generator for work.
func NewWork(title, description, author string, duration int, tags []string, when time.Time) *Work {
	nowString := time.Now().Format(time.RFC3339)
	now, err := time.Parse(time.RFC3339, nowString)
	if err != nil {
		fmt.Println("now is not in a valid time format.")
		os.Exit(1)
	}
	sort.Strings(tags)
	return &Work{
		Title:       title,
		Description: description,
		Author:      author,
		Where:       "",
		Duration:    duration,
		Tags:        tags,
		When:        when,
		Created:     now,
	}
}

func (w Work) String() string {
	finalString := " "
	if w.Title != "" {
		finalString = fmt.Sprintf("%s Title: %s,", finalString, w.Title)
	}
	if w.Description != "" {
		finalString = fmt.Sprintf("%s Description: %s,", finalString, w.Description)
	}
	if w.Author != "" {
		finalString = fmt.Sprintf("%s Author: %s,", finalString, w.Author)
	}
	if w.Where != "" {
		finalString = fmt.Sprintf("%s Where: %s,", finalString, w.Where)
	}
	if w.Duration != 0 {
		finalString = fmt.Sprintf("%s Duration: %d,", finalString, w.Duration)
	}
	if len(w.Tags) > 0 {
		finalString = fmt.Sprintf("%s Tags: %s,", finalString, strings.Join(w.Tags, ", "))
	}
	if !w.When.Equal(time.Time{}) {
		finalString = fmt.Sprintf("%s When: %s,", finalString, w.When)
	}
	if !w.Created.Equal(time.Time{}) {
		finalString = fmt.Sprintf("%s Created: %s,", finalString, w.Created)
	}
	return strings.TrimSpace(finalString[:len(finalString)-1])
}
