package integration

import (
	"fmt"
	"io/ioutil"
	"sort"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
	"gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
)

const binaryName = "worklog"

const length = 56

var now = time.Now()
var tmUTC = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

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
	wls, _ := ymlRepo.GetAllBetweenDates(tmUTC.Add(time.Hour*1*-1), tmUTC.Add(time.Hour*1), exp)

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
