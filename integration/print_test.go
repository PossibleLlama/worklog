package integration

import (
	"os"
	"os/exec"
	"path"
	"testing"

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

	var tests = []struct {
		name      string
		args      []string
		success   bool
		expOutput string
	}{
		{
			name:      "No arguments",
			args:      []string{},
			success:   false,
			expOutput: printNoArgs,
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

			assert.Equal(t, testItem.expOutput, string(output))
			if testItem.success {
				assert.Nil(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
