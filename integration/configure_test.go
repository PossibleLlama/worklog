package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/assert"
)

const overrideNoArgs = `Error: defaults requires at least one argument
Usage:
  worklog configure overrideDefaults [flags]

Flags:
      --author string   The authour for all work
      --duration int    Default duration that work takes (default -1)
      --format string   Format to print work in. If provided, must be one of 'pretty', 'yaml', 'json'
  -h, --help            help for overrideDefaults

Global Flags:
      --config string   config file (default is $HOME/.worklog/config.yml)

defaults requires at least one argument
`
const overrideInvalidFormat = `Error: format is not valid
Usage:
  worklog configure overrideDefaults [flags]

Flags:
      --author string   The authour for all work
      --duration int    Default duration that work takes (default -1)
      --format string   Format to print work in. If provided, must be one of 'pretty', 'yaml', 'json'
  -h, --help            help for overrideDefaults

Global Flags:
      --config string   config file (default is $HOME/.worklog/config.yml)

format is not valid
`

func TestConfigure(t *testing.T) {
	randString := helpers.RandString(length)

	var tests = []struct {
		name      string
		args      []string
		success   bool
		expOutput string
		expFile   *model.Config
	}{
		{
			name:      "Defaults are used",
			args:      []string{},
			success:   true,
			expOutput: "Successfully configured\n",
			expFile: &model.Config{
				Author: "",
				Defaults: model.Defaults{
					Duration: 15,
					Format:   "pretty",
				},
			},
		}, {
			name:      "Override defaults requires an argument",
			args:      []string{"overrideDefaults"},
			success:   false,
			expOutput: overrideNoArgs,
			expFile:   nil,
		}, {
			name:      "Override defaults has author",
			args:      []string{"overrideDefaults", "--author", randString},
			success:   true,
			expOutput: "Successfully configured\n",
			expFile: &model.Config{
				Author: randString,
				Defaults: model.Defaults{
					Duration: 15,
					Format:   "",
				},
			},
		}, {
			name: "Override defaults has duration",
			args: []string{"overrideDefaults",
				"--duration",
				fmt.Sprintf("%d", length)},
			success:   true,
			expOutput: "Successfully configured\n",
			expFile: &model.Config{
				Author: "",
				Defaults: model.Defaults{
					Duration: length,
					Format:   "",
				},
			},
		}, {
			name:      "Override defaults has format pretty",
			args:      []string{"overrideDefaults", "--format", "pretty"},
			success:   true,
			expOutput: "Successfully configured\n",
			expFile: &model.Config{
				Author: "",
				Defaults: model.Defaults{
					Duration: 15,
					Format:   "pretty",
				},
			},
		}, {
			name:      "Override defaults has format yaml",
			args:      []string{"overrideDefaults", "--format", "yaml"},
			success:   true,
			expOutput: "Successfully configured\n",
			expFile: &model.Config{
				Author: "",
				Defaults: model.Defaults{
					Duration: 15,
					Format:   "yaml",
				},
			},
		}, {
			name:      "Override defaults has format json",
			args:      []string{"overrideDefaults", "--format", "json"},
			success:   true,
			expOutput: "Successfully configured\n",
			expFile: &model.Config{
				Author: "",
				Defaults: model.Defaults{
					Duration: 15,
					Format:   "json",
				},
			},
		}, {
			name:      "Override defaults has random format",
			args:      []string{"overrideDefaults", "--format", randString},
			success:   false,
			expOutput: overrideInvalidFormat,
			expFile:   nil,
		}, {
			name:      "Override defaults has multiple args",
			args:      []string{"overrideDefaults", "--author", randString, "--format", "json"},
			success:   true,
			expOutput: "Successfully configured\n",
			expFile: &model.Config{
				Author: randString,
				Defaults: model.Defaults{
					Duration: 15,
					Format:   "json",
				},
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			dir, err := os.Getwd()
			if err != nil {
				t.Error(err)
			}

			testItem.args = append([]string{"configure"}, testItem.args...)
			cmd := exec.Command(path.Join(dir, binaryName), testItem.args...)
			output, err := cmd.CombinedOutput()

			assert.Equal(t, testItem.expOutput, string(output))
			if testItem.success {
				assert.Nil(t, err)
				assert.Equal(t, testItem.expFile, getActualConfig(t))
			} else {
				assert.Error(t, err)
			}
		})
	}
}
