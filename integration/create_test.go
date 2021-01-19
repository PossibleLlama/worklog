package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	randString := helpers.RandString(length)
	tm := time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		time.Now().Hour(),
		time.Now().Minute(),
		time.Now().Second(),
		0,
		time.UTC)

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
			expOutput: fmt.Sprintf("Error: required flag(s) \"title\" not set\nUsage:\n  worklog create [flags]\n\nFlags:\n      --description string   A description of the work\n      --duration int         Length of time spent on the work (default -1)\n  -h, --help                 help for create\n      --tags string          Comma seperated list of tags this work relates to\n      --title string         A short description of the work done\n      --when string          When the work was worked in RFC3339 format (default \"%s\")\n\nGlobal Flags:\n      --config string   config file (default is $HOME/.worklog/config.yml)\n\nrequired flag(s) \"title\" not set\n", time.Now().Format(time.RFC3339)),
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
			args:      []string{"--title", "Create with when", "--when", tm.Format(time.RFC3339)},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title: "Create with when",
				When:  tm,
			},
		}, {
			name:      "Create with all",
			args:      []string{"--title", "Create with all", "--description", randString, "--duration", fmt.Sprintf("%d", length), "--tags", randString},
			success:   true,
			expOutput: "Saving file...\nSaved file\n",
			expFile: &model.Work{
				Title:       "Create with all",
				Description: randString,
				Tags:        []string{randString},
				When:        tm,
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
			output, err := cmd.CombinedOutput()

			assert.Equal(t, testItem.expOutput, string(output))
			if testItem.success {
				actualFile := getActualWork(t, fmt.Sprintf("%d-%02d-%02d_%02d:%02d_%s.yml",
					tm.Year(),
					int(tm.Month()),
					tm.Day(),
					tm.Hour(),
					tm.Minute(),
					strings.ReplaceAll(testItem.name, " ", "_")))
				cfg := getActualConfig(t)

				assert.NotEmpty(t, actualFile.ID)
				assert.Len(t, actualFile.ID, 16)
				assert.Equal(t, testItem.expFile.Title, actualFile.Title)
				assert.NotEqual(t, time.Time{}, actualFile.CreatedAt)
				assert.Equal(t, "", actualFile.Where)
				assert.Equal(t, 1, actualFile.Revision)

				if testItem.expFile.Description != "" {
					assert.Equal(t, testItem.expFile.Description, actualFile.Description)
				}

				if testItem.expFile.Duration > 0 {
					assert.Equal(t, testItem.expFile.Duration, actualFile.Duration)
				} else {
					if cfg.Defaults.Duration > 0 {
						assert.Equal(t, cfg.Defaults.Duration, actualFile.Duration)
					} else {
						assert.Equal(t, 0, actualFile.Duration)
					}
				}

				if testItem.expFile.Author != "" {
					assert.Equal(t, testItem.expFile.Author, actualFile.Author)
				} else {
					if cfg.Author != "" {
						assert.Equal(t, cfg.Author, actualFile.Author)
					} else {
						assert.Equal(t, "", actualFile.Author)
					}
				}

				if !testItem.expFile.When.Equal(time.Time{}) {
					assert.Equal(t, testItem.expFile.When, actualFile.When)
				} else {
					assert.NotEqual(t, time.Time{}, actualFile.When)
				}

				if len(testItem.expFile.Tags) != 0 {
					assert.Equal(t, testItem.expFile.Tags, actualFile.Tags)
				}
			}
		})
	}
}
