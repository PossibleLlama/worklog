package main

import (
	"fmt"
	"os"
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
	When        time.Time
	Created     time.Time
}

// New is the generator for work.
func New(title, description, author, where, when string) *Work {
	nowString := time.Now().Format(time.RFC3339)
	now, err := time.Parse(time.RFC3339, nowString)
	if err != nil {
		fmt.Println("now is not in a valid time format.")
		os.Exit(1)
	}
	if len(when) == 0 {
		when = nowString
	} else if len(when) == 10 {
		when = fmt.Sprintf("%sT00:00:00Z", when)
	}
	whenAsDate, err := time.Parse(time.RFC3339, when)
	if err != nil {
		fmt.Printf("when is not in a valid time format. %s\n", err.Error())
		os.Exit(1)
	}
	return &Work{
		Title:       title,
		Description: description,
		Author:      author,
		Where:       where,
		When:        whenAsDate,
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
	if ! w.When.Equal(time.Time{}) {
		finalString = fmt.Sprintf("%s When: %s,", finalString, w.When)
	}
	if ! w.Created.Equal(time.Time{}) {
		finalString = fmt.Sprintf("%s Created: %s,", finalString, w.Created)
	}
	return strings.TrimSpace(finalString[:len(finalString)-1])
}