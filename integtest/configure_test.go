package integtest

import (
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	var tests = []struct {
		name      string
		args      []string
		expOutput string
		expFile   model.Config
	}{
		{
			name:      "defaults are used",
			args:      []string{},
			expOutput: "Successfully configured\n",
			expFile: model.Config{
				Author: "",
				Defaults: model.Defaults{
					Duration: 15,
					Format:   "pretty",
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
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, testItem.expOutput, string(output))
			assert.Equal(t, testItem.expFile, *getActualConfig(t))
		})
	}
}
