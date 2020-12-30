package model

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
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
	if w.Duration != 0 {
		finalString = fmt.Sprintf("%s Duration: %d,", finalString, w.Duration)
	}
	if len(w.Tags) > 0 {
		finalString = fmt.Sprintf("%s Tags: [%s],", finalString, strings.Join(w.Tags, ", "))
	}
	if !w.When.Equal(time.Time{}) {
		finalString = fmt.Sprintf("%s When: %s,", finalString, helpers.TimeFormat(w.When))
	}
	if !w.CreatedAt.Equal(time.Time{}) {
		finalString = fmt.Sprintf("%s CreatedAt: %s", finalString, helpers.TimeFormat(w.CreatedAt))
	}
	return strings.TrimSpace(finalString[:len(finalString)-1])
}

// PrettyString works like string, but with greater spacing and line breaks
func (w Work) PrettyString() string {
	pw := workToPrintWork(w)
	finalString := " "
	if pw.Title != "" {
		finalString = fmt.Sprintf("%sTitle: %s\n", finalString, pw.Title)
	}
	if pw.Description != "" {
		finalString = fmt.Sprintf("%sDescription: %s\n", finalString, pw.Description)
	}
	if pw.Author != "" {
		finalString = fmt.Sprintf("%sAuthor: %s\n", finalString, pw.Author)
	}
	if pw.Duration != 0 {
		finalString = fmt.Sprintf("%sDuration: %d\n", finalString, pw.Duration)
	}
	if len(pw.Tags) > 0 {
		finalString = fmt.Sprintf("%sTags: [%s]\n", finalString, strings.Join(pw.Tags, ", "))
	}
	if !pw.When.Equal(time.Time{}) {
		finalString = fmt.Sprintf("%sWhen: %s\n", finalString, helpers.TimeFormat(pw.When))
	}
	return strings.TrimSpace(finalString[:len(finalString)-1])
}

// WriteText takes a writer and outputs a text representation of Work to it
func (w Work) WriteText(writer io.Writer) error {
	_, err := writer.Write([]byte(w.PrettyString()))
	return err
}

// WriteAllWorkToText takes a writer and list of work, and outputs a text representation of Work to the writer
func WriteAllWorkToText(writer io.Writer, w []*Work) error {
	for index, work := range w {
		err := work.WriteText(os.Stdout)
		if err != nil {
			return err
		}
		if index != len(w)-1 {
			fmt.Println()
		}
		fmt.Println()
	}
	return nil
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

// WriteAllWorkToYAML takes a writer and list of work, and outputs a YAML representation of Work to the writer
func WriteAllWorkToYAML(writer io.Writer, w []*Work) error {
	pw := []printWork{}
	for _, work := range w {
		pw = append(pw, workToPrintWork(*work))
	}

	b, err := yaml.Marshal(pw)
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

// WriteAllWorkToJSON takes a writer and list of work, and outputs a JSON representation of Work to the writer
func WriteAllWorkToJSON(writer io.Writer, w []*Work) error {
	pw := []printWork{}
	for _, work := range w {
		pw = append(pw, workToPrintWork(*work))
	}

	b, err := json.Marshal(pw)
	if err != nil {
		return err
	}

	_, err = writer.Write(b)
	return err
}
