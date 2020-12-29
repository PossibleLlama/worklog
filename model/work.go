package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// Work data model.
// stores information as to what
// work was done by who and when.
type Work struct {
	Title       string    `json:"title" yaml:"title"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Author      string    `json:"author,omitempty" yaml:"author,omitempty"`
	Where       string    `json:"where,omitempty" yaml:"where,omitempty"`
	Duration    int       `json:"duration" yaml:"duration"`
	Tags        []string  `json:"tags,flow,omitempty" yaml:"tags,flow,omitempty"`
	When        time.Time `json:"when" yaml:"when"`
	CreatedAt   time.Time `json:"createdAt" yaml:"createdAt"`
}

type printWork struct {
	Title       string    `json:"title" yaml:"title"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Author      string    `json:"author,omitempty" yaml:"author,omitempty"`
	Duration    int       `json:"duration" yaml:"duration"`
	Tags        []string  `json:"tags,flow,omitempty" yaml:"tags,flow,omitempty"`
	When        time.Time `json:"when" yaml:"when"`
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
		CreatedAt:   now,
	}
}

func workToPrintWork(w Work) printWork {
	return printWork{
		Title:       w.Title,
		Description: w.Description,
		Author:      w.Author,
		Duration:    w.Duration,
		Tags:        w.Tags,
		When:        w.When,
	}
}

// String generates a stringified version of the Work
func (w Work) String() string {
	pw := workToPrintWork(w)
	finalString := " "
	if pw.Title != "" {
		finalString = fmt.Sprintf("%s Title: %s,", finalString, pw.Title)
	}
	if pw.Description != "" {
		finalString = fmt.Sprintf("%s Description: %s,", finalString, pw.Description)
	}
	if pw.Author != "" {
		finalString = fmt.Sprintf("%s Author: %s,", finalString, pw.Author)
	}
	if pw.Duration != 0 {
		finalString = fmt.Sprintf("%s Duration: %d,", finalString, pw.Duration)
	}
	if len(pw.Tags) > 0 {
		finalString = fmt.Sprintf("%s Tags: [%s],", finalString, strings.Join(pw.Tags, ", "))
	}
	if !pw.When.Equal(time.Time{}) {
		finalString = fmt.Sprintf("%s When: %s,", finalString, pw.When)
	}
	return strings.TrimSpace(finalString[:len(finalString)-1])
}

// WriteText takes a writer and outputs a text representation of Work to it
func (w Work) WriteText(writer io.Writer) error {
	_, err := writer.Write([]byte(w.String()))
	return err
}

// WriteYAML takes a writer and outputs a YAML representation of Work to it
func (w Work) WriteYAML(writer io.Writer) error {
	b, err := yaml.Marshal(workToPrintWork(w))
	if err != nil {
		return err
	}

	_, err = writer.Write(b)
	return err
}

// WriteJSON takes a writer and outputs a JSON representation of Work to it
func (w Work) WriteJSON(writer io.Writer) error {
	b, err := json.Marshal(workToPrintWork(w))
	if err != nil {
		return err
	}

	_, err = writer.Write(b)
	return err
}
