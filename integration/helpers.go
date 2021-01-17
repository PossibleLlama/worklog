package integration

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/PossibleLlama/worklog/model"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
)

const binaryName = "worklog"

const length = 56

func getActualConfig(t *testing.T) *model.Config {
	var actualFile model.Config
	home, err := homedir.Dir()
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/.worklog/config.yml", home))
	if err != nil {
		t.Error(err)
	}

	err = yaml.Unmarshal(file, &actualFile)
	if err != nil {
		t.Error(err)
	}

	return &actualFile
}
