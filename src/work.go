package main

import (
	"fmt"
	"os"
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
	whenAsDate, err := time.Parse(time.RFC3339, when)
	if err != nil {
		fmt.Printf("when is not in a valid time format.")
		os.Exit(1)
	}
	return &Work{
		Title:       title,
		Description: description,
		Author:      author,
		Where:       where,
		When:        whenAsDate,
		Created:     time.Now(),
	}
}
