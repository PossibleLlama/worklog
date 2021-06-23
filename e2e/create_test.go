package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/assert"
)

var createNoArgs = fmt.Sprintf(`Error: required flag(s) "title" not set
Usage:
  worklog create [flags]

Flags:
      --author string        The author of the work
      --description string   A description of the work
      --duration int         Length of time spent on the work (default -1)
  -h, --help                 help for create
      --tags string          Comma separated list of tags this work relates to
      --title string         A short description of the work done
      --when string          When the work was worked in RFC3339 format (default "%s")

Global Flags:
      --config string   config file including file extension (default ".worklog/config.yml")
      --legacy          Use legacy yaml repository for storing/retrieving worklogs
      --repo string     repository that worklogs are stored in (default ".worklog/worklog.db")

required flag(s) "title" not set
`, time.Now().Format(time.RFC3339))

func TestCreate(t *testing.T) {
	randString := helpers.RandAlphabeticString(length)

	var tests = []struct {
		name      string
		args      []string
		success   bool
		expOutput string
		expFile   *model.Work
	}{
		{
			name:      "Create new",
			args:      []string{"--title", "Create new"},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title: "Create new",
			},
		}, {
			name:      "No arguments",
			args:      []string{},
			success:   false,
			expOutput: createNoArgs,
			expFile:   nil,
		}, {
			name:      "Create with description",
			args:      []string{"--title", "Create with description", "--description", randString},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title:       "Create with description",
				Description: randString,
			},
		}, {
			name:      "Create with author",
			args:      []string{"--title", "Create with author", "--author", randString},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title:  "Create with author",
				Author: randString,
			},
		}, {
			name:      "Create with duration",
			args:      []string{"--title", "Create with duration", "--duration", fmt.Sprintf("%d", length)},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title:    "Create with duration",
				Duration: length,
			},
		}, {
			name:      "Create with tags",
			args:      []string{"--title", "Create with tags", "--tags", randString},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title: "Create with tags",
				Tags:  []string{randString},
			},
		}, {
			name:      "Create with when",
			args:      []string{"--title", "Create with when", "--when", tmUTC.Format(time.RFC3339)},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title: "Create with when",
				When:  tmUTC,
			},
		}, {
			name:      "Create with all",
			args:      []string{"--title", "Create with all", "--description", randString, "--duration", fmt.Sprintf("%d", length), "--tags", randString, "--author", randString},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title:       "Create with all",
				Description: randString,
				Author:      randString,
				Tags:        []string{randString},
				When:        tmUTC,
				Duration:    length,
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Error(err)
			}

			testItem.args = append([]string{"create"}, testItem.args...)
			cmd := exec.Command(path.Join(dir, binaryName), testItem.args...)
			output, _ := cmd.CombinedOutput()

			assert.Equal(t, testItem.expOutput, string(output))

			cfg := getActualConfig(t)
			if testItem.success {
				actualFile := getActualWork(t, testItem.expFile, cfg)

				assert.NotEmpty(t, actualFile.ID, "ID is empty")
				assert.Len(t, actualFile.ID, 20, "ID is not of length 20")
				assert.Equal(t, testItem.expFile.Title, actualFile.Title, "Title does not match")
				assert.NotEqual(t, time.Time{}, actualFile.CreatedAt, "CreatedAt has defaulted")
				assert.Equal(t, 1, actualFile.Revision, "Revision is not 1")

				if testItem.expFile.Description != "" {
					assert.Equal(t, testItem.expFile.Description, actualFile.Description, "Description does not match")
				}

				if testItem.expFile.Duration > 0 {
					assert.Equal(t, testItem.expFile.Duration, actualFile.Duration, "Duration does not match provided")
				} else {
					if cfg.Defaults.Duration > 0 {
						assert.Equal(t, cfg.Defaults.Duration, actualFile.Duration, "Description does not match config")
					} else {
						assert.Equal(t, 0, actualFile.Duration, "Duration does not match default")
					}
				}

				if testItem.expFile.Author != "" {
					assert.Equal(t, testItem.expFile.Author, actualFile.Author, "Author does not match provided")
				} else {
					if cfg.Defaults.Author != "" {
						assert.Equal(t, cfg.Defaults.Author, actualFile.Author, "Author does not match config")
					} else {
						assert.Equal(t, "", actualFile.Author, "Author does not match default")
					}
				}

				if !testItem.expFile.When.Equal(time.Time{}) {
					assert.Equal(t, testItem.expFile.When.Year(), actualFile.When.Year(), "When's year does not match provided")
					assert.Equal(t, testItem.expFile.When.Month(), actualFile.When.Month(), "When's month does not match provided")
					assert.Equal(t, testItem.expFile.When.Day(), actualFile.When.Day(), "When's day does not match provided")
					assert.Equal(t, testItem.expFile.When.Hour(), actualFile.When.Hour(), "When's hour does not match provided")
					assert.Equal(t, testItem.expFile.When.Minute(), actualFile.When.Minute(), "When's minute does not match provided")
				} else {
					assert.NotEqual(t, time.Time{}, actualFile.When, "When didn't default")
				}

				if len(testItem.expFile.Tags) != 0 {
					assert.Equal(t, testItem.expFile.Tags, actualFile.Tags, "Tags do not match")
				}
			}
		})
	}
}
