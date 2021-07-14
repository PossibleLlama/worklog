package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
	"gopkg.in/yaml.v2"
)

const (
	binaryName = "worklog"
	repoName   = "e2e.db"
	configName = "e2e.yml"
)

const length = 56

var now = time.Now()
var tmUTC = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

func execBinary(args ...string) (string, error) {
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		panic(dirErr)
	}
	cmd := exec.Command(
		path.Join(dir, binaryName),
		args...)
	output, cmdErr := cmd.CombinedOutput()
	return string(output), cmdErr
}

func execConfiguredBinary(args ...string) (string, error) {
	dir, dirErr := os.Getwd()
	if dirErr != nil {
		panic(dirErr)
	}
	args = append([]string{
		//"--repo \"local\"",
		//fmt.Sprintf("--repoPath \"%s\"", path.Join(dir, repoName)),
		fmt.Sprintf("--config \"%s\"", path.Join(dir, configName)),
	}, args...)
	return execBinary(args...)
}

func getActualConfig(t *testing.T) *model.Config {
	var actualFile model.Config
	dir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	file, err := ioutil.ReadFile(
		path.Join(dir, configName))
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

	dir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	_ = repository.NewYamlConfig(path.Join(dir, configName))
	repo := repository.NewBBoltRepo(path.Join(dir, repoName))
	wls, _ := repo.GetAllBetweenDates(tmUTC.Add(time.Hour*1*-1), tmUTC.Add(time.Hour*1), exp)

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
