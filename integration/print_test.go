package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const printNoArgs = `Error: one flag is required
Usage:
  worklog print [flags]

Flags:
      --author string        Filter by work including author
      --description string   Filter by work including description
      --endDate string       Date till which to find worklogs. Only functions in conjunction with startDate
  -h, --help                 help for print
      --json                 Output in a json format
      --pretty               Output in a text format
      --startDate string     Date from which to find worklogs
      --tags string          Filter by work including all tags
      --thisWeek             Prints this weeks work
      --title string         Filter by work including title
      --today                Print today's work
      --yaml                 Output in a yaml format

Global Flags:
      --config string   config file (default is $HOME/.worklog/config.yml)

one flag is required
`

func TestPrint(t *testing.T) {
	tm := time.Now().Add(time.Hour * length * length * -1)
	var tests = []struct {
		name      string
		args      []string
		success   bool
		expOutput string
	}{
		{
			name:      "Print with start and end dates pretty, no wl",
			args:      []string{"--startDate", tm.Format(time.RFC3339), "--endDate", tm.Format(time.RFC3339), "--pretty"},
			success:   true,
			expOutput: fmt.Sprintf("No work found between %02d-%02d-%02d 00:00:00 +0000 UTC and %02d-%02d-%02d 23:59:59 +0000 UTC with the given filter\n", tm.Year(), int(tm.Month()), tm.Day(), tm.Year(), int(tm.Month()), tm.Day()),
		}, {
			name:      "No arguments",
			args:      []string{},
			success:   false,
			expOutput: printNoArgs,
		},
		// Relies on the create_test.go tests having been ran,
		// and the wl's generated from that.
		{
			name:      "Print with start and end dates pretty, multiple wl",
			args:      []string{"--startDate", time.Now().Format(time.RFC3339), "--endDate", time.Now().Format(time.RFC3339), "--pretty"},
			success:   true,
			expOutput: fmt.Sprintf("\nTitle: Create new\nAuthor: %s\nDuration: 15\nTags: []\nWhen: ", getActualConfig(t).Defaults.Author),
		}, {
			name:      "Print with start and end dates yaml, multiple wl",
			args:      []string{"--startDate", time.Now().Format(time.RFC3339), "--endDate", time.Now().Format(time.RFC3339), "--yaml"},
			success:   true,
			expOutput: fmt.Sprintf("\n  title: Create new\n  author: %s\n  duration: 15\n  tags: [\"\"]\n  when: ", getActualConfig(t).Defaults.Author),
		}, {
			name:      "Print with start and end dates json, multiple wl",
			args:      []string{"--startDate", time.Now().Format(time.RFC3339), "--endDate", time.Now().Format(time.RFC3339), "--json"},
			success:   true,
			expOutput: fmt.Sprintf("\",\"title\":\"Create new\",\"author\":\"%s\",\"duration\":15,\"tags\":[\"\"],\"when\":", getActualConfig(t).Defaults.Author),
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