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

var createNoArgs = fmt.Sprintf(`Error: required flag(s) "title" not set
Usage:
  worklog create [flags]

Flags:
      --description string   A description of the work
      --duration int         Length of time spent on the work (default -1)
  -h, --help                 help for create
      --tags string          Comma seperated list of tags this work relates to
      --title string         A short description of the work done
      --when string          When the work was worked in RFC3339 format (default "%s")

Global Flags:
      --config string   config file (default is $HOME/.worklog/config.yml)

required flag(s) "title" not set
`, time.Now().Format(time.RFC3339))

func TestCreate(t *testing.T) {
	var tests = []struct {
		name      string
		args      []string
		expOutput string
	}{
		{
			name:      "Create new",
			args:      []string{"--title", helpers.RandString(length)},
			expOutput: "Saving file...\nSaved file\n",
		}, {
			name:      "No arguments",
			args:      []string{},
			expOutput: createNoArgs,
		}, {
			name:      "Create with description",
			args:      []string{"--title", helpers.RandString(length), "--description", helpers.RandString(length)},
			expOutput: "Saving file...\nSaved file\n",
		}, {
			name:      "Create with duration",
			args:      []string{"--title", helpers.RandString(length), "--duration", fmt.Sprintf("%d", length)},
			expOutput: "Saving file...\nSaved file\n",
		}, {
			name:      "Create with tags",
			args:      []string{"--title", helpers.RandString(length), "--tags", helpers.RandString(length)},
			expOutput: "Saving file...\nSaved file\n",
		}, {
			name:      "Create with when",
			args:      []string{"--title", helpers.RandString(length), "--when", time.Now().Format(time.RFC3339)},
			expOutput: "Saving file...\nSaved file\n",
		}, {
			name:      "Create with all",
			args:      []string{"--title", helpers.RandString(length), "--description", helpers.RandString(length), "--duration", fmt.Sprintf("%d", length), "--tags", helpers.RandString(length)},
			expOutput: "Saving file...\nSaved file\n",
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
		})
	}
}
