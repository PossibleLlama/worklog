package integration

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/model"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
)

const binaryName = "worklog"

const length = 56

// If we are in the first half of the second.
// Will hopefully cut out the failures when a second rolls over
// and causes values to not match.
func TestMain(m *testing.M) {
	ran := false
	for i := 0; i < 10; i++ {
		now := time.Now()
		if now.Second() == now.Round(time.Second).Second() {
			ran = true
			m.Run()
		}
		time.Sleep(time.Millisecond * 100)
	}
	if !ran {
		fmt.Print("Didn't pause correctly")
		m.Run()
	}
}

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

func getActualWork(t *testing.T, name string) *model.Work {
	var actualFile model.Work
	home, err := homedir.Dir()
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/.worklog/%s", home, name))
	if err != nil {
		t.Error(err)
	}

	err = yaml.Unmarshal(file, &actualFile)
	if err != nil {
		t.Error(err)
	}

	return &actualFile
}
