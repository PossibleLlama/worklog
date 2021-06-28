package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/stretchr/testify/assert"
)

const rootHelp = `A CLI tool to let people track what work they
have completed. You can record what work you do,
and get a summary of what you've done each day.

For information on using the CLI, use worklog
--help

Usage:
  worklog [command]

Available Commands:
  configure   Setup configuration for the application
  create      Create a new record of work
  edit        Edit an existing record of work
  help        Help about any command
  print       Print all worklogs since provided date

Flags:
      --config string     Config file including file extension (default ".worklog/config.yml")
  -h, --help              help for worklog
      --repo string       Which type of repository to use for storing/retrieving worklogs (default "legacy")
      --repoPath string   Directory path for repository that worklogs are stored in (default ".worklog/worklog.db")
  -v, --version           version for worklog

Use "worklog [command] --help" for more information about a command.
`

func TestRoot(t *testing.T) {
	var tests = []struct {
		name      string
		args      []string
		expOutput string
	}{
		{
			name:      "No arguments",
			args:      []string{},
			expOutput: rootHelp,
		}, {
			name:      "Help",
			args:      []string{"--help"},
			expOutput: rootHelp,
		}, {
			name:      "Version",
			args:      []string{"--version"},
			expOutput: fmt.Sprintf("worklog version %s\n", helpers.Version),
		}, {
			name:      "Change config file",
			args:      []string{"--config", "/worklog-test/config.yml"},
			expOutput: rootHelp,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Error(err)
			}

			cmd := exec.Command(path.Join(dir, binaryName), testItem.args...)
			output, err := cmd.CombinedOutput()

			assert.Equal(t, testItem.expOutput, string(output))
			assert.Nil(t, err)
		})
	}
}
