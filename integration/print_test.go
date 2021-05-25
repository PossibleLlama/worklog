package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/stretchr/testify/assert"
)

const printNoArgs = `Error: one flag is required
Usage:
  worklog print [flags]

Flags:
  -a, --all                  Output all fields of the worklog
      --author string        Filter by work including author
      --description string   Filter by work including description
      --endDate string       Date till which to find worklogs. Only functions in conjunction with startDate
  -h, --help                 help for print
  -j, --json                 Output in a json format
  -p, --pretty               Output in a text format
      --startDate string     Date from which to find worklogs
      --tags string          Filter by work including all tags
  -w, --thisWeek             Prints this weeks work
      --title string         Filter by work including title
  -t, --today                Print today's work
  -y, --yaml                 Output in a yaml format

Global Flags:
      --config string   config file (default is $HOME/.worklog/config.yml)

one flag is required
`

func TestPrint(t *testing.T) {
	unusedDate := time.Date(1980, time.February, 10, 1, 2, 3, 0, time.Now().Location())

	randStr := helpers.RandAlphabeticString(20)
	var tests = []struct {
		name      string
		args      []string
		success   bool
		expOutput string
	}{
		{
			name:      "Print with start and end dates pretty, no wl",
			args:      []string{"--startDate", unusedDate.Format(time.RFC3339), "--endDate", unusedDate.Format(time.RFC3339), "--pretty"},
			success:   true,
			expOutput: fmt.Sprintf("No work found between %02d-%02d-%02d 00:00:00 +0000 UTC and %02d-%02d-%02d 23:59:59 +0000 UTC with the given filter\n", unusedDate.Year(), int(unusedDate.Month()), unusedDate.Day(), unusedDate.Year(), int(unusedDate.Month()), unusedDate.Day()),
		}, {
			name:      "Print with invalid ID pretty",
			args:      []string{randStr, "--pretty"},
			success:   true,
			expOutput: fmt.Sprintf("No work found between 0001-01-01 00:00:00 +0000 UTC and 0000-12-31 23:59:59 +0000 UTC with the given filter with id's [%s]", randStr),
		}, {
			name:      "Print with invalid ID json",
			args:      []string{randStr, "--json"},
			success:   true,
			expOutput: "[]",
		}, {
			name:      "Print with valid ID pretty",
			args:      []string{"a", "--pretty"},
			success:   false,
			expOutput: "Error: ID 'a' is not unique",
		},
		// TODO give certainty to what ID's exist
		{
			name:      "No arguments",
			args:      []string{},
			success:   false,
			expOutput: printNoArgs,
		},
		// Relies on the create_test.go tests having been ran,
		// and the wl's generated from that.
		{
			name:      "Print all fields",
			args:      []string{"--startDate", tmUTC.Format(time.RFC3339), "--endDate", tmUTC.Format(time.RFC3339), "--all", "--pretty"},
			success:   true,
			expOutput: fmt.Sprintf("\nRevision: 1\nTitle: Create new\nAuthor: %s\nDuration: %d\nWhen: ", getActualConfig(t).Defaults.Author, getActualConfig(t).Defaults.Duration),
		}, {
			name:      "Print all fields, shorthand",
			args:      []string{"--startDate", tmUTC.Format(time.RFC3339), "--endDate", tmUTC.Format(time.RFC3339), "-ap"},
			success:   true,
			expOutput: fmt.Sprintf("\nRevision: 1\nTitle: Create new\nAuthor: %s\nDuration: %d\nWhen: ", getActualConfig(t).Defaults.Author, getActualConfig(t).Defaults.Duration),
		}, {
			name:      "Print with start and end dates pretty, multiple wl",
			args:      []string{"--startDate", tmUTC.Format(time.RFC3339), "--endDate", tmUTC.Format(time.RFC3339), "--pretty"},
			success:   true,
			expOutput: fmt.Sprintf("\nTitle: Create new\nAuthor: %s\nDuration: %d\nWhen: ", getActualConfig(t).Defaults.Author, getActualConfig(t).Defaults.Duration),
		}, {
			name:      "Print with start and end dates yaml, multiple wl",
			args:      []string{"--startDate", tmUTC.Format(time.RFC3339), "--endDate", tmUTC.Format(time.RFC3339), "--yaml"},
			success:   true,
			expOutput: fmt.Sprintf("\n  title: Create new\n  author: %s\n  duration: %d\n  when: ", getActualConfig(t).Defaults.Author, getActualConfig(t).Defaults.Duration),
		}, {
			name:      "Print with start and end dates json, multiple wl",
			args:      []string{"--startDate", tmUTC.Format(time.RFC3339), "--endDate", tmUTC.Format(time.RFC3339), "--json"},
			success:   true,
			expOutput: fmt.Sprintf("\",\"title\":\"Create new\",\"author\":\"%s\",\"duration\":%d,\"when\":", getActualConfig(t).Defaults.Author, getActualConfig(t).Defaults.Duration),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Error(err)
			}

			testItem.args = append([]string{"print"}, testItem.args...)
			cmd := exec.Command(path.Join(dir, binaryName), testItem.args...)
			output, err := cmd.CombinedOutput()

			assert.Contains(t, string(output), testItem.expOutput)
			if testItem.success {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
