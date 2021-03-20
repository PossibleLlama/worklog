package integration

import (
	"fmt"
	"io/ioutil"
	"sort"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
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
		if time.Now().Nanosecond() < 50 {
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

func getActualWork(t *testing.T, exp *model.Work, cfg *model.Config) *model.Work {
	if exp.Author == "" {
		exp.Author = cfg.Defaults.Author
	}
	ymlRepo := repository.NewYamlFileRepo()
	wls, _ := ymlRepo.GetAllBetweenDates(helpers.Midnight(time.Now()), helpers.Midnight(time.Now()).Add(time.Hour*24), exp)

	fmt.Printf("\n********\n%d WL's found for filter %s\n", len(wls), exp)

	var actual *model.Work
	switch len(wls) {
	case 0:
		t.Errorf("Unable to find work with filter %s", exp)
		t.FailNow()
	case 1:
		actual = wls[0]
	default:
		fmt.Printf("Items\n%s\n", wls)
		dateSortedWls := make(model.WorkList, 0, len(wls))
		for _, d := range wls {
			dateSortedWls = append(dateSortedWls, d)
		}
		sort.Sort(dateSortedWls)
		actual = wls[len(wls)-1]
	}
	return actual
}
