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
			expOutput: fmt.Sprintf("Error: required flag(s) \"title\" not set\nUsage:\n  worklog create [flags]\n\nFlags:\n      --description string   A description of the work\n      --duration int         Length of time spent on the work (default -1)\n  -h, --help                 help for create\n      --tags string          Comma seperated list of tags this work relates to\n      --title string         A short description of the work done\n      --when string          When the work was worked in RFC3339 format (default \"%s\")\n\nGlobal Flags:\n      --config string   config file (default is $HOME/.worklog/config.yml)\n\nrequired flag(s) \"title\" not set\n", time.Now().Format(time.RFC3339)),
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
