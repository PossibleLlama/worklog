package integtest

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/PossibleLlama/worklog/model"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
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

			var actualFile model.Config
			home, err := homedir.Dir()
			if err != nil {
				t.Error(err)
			}
			file, err := ioutil.ReadFile(fmt.Sprintf("%s/.worklog/config.yml", home))
			if err != nil {
				t.Error(err)
			}

			actualOutput := string(output)
			err = yaml.Unmarshal(file, &actualFile)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, testItem.expOutput, actualOutput)
			assert.Equal(t, testItem.expFile, actualFile)
		})
	}
}
