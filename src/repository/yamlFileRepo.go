package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

type yamlFileRepo struct{}

// NewYamlFileRepo Generator for repository storing worklogs
// on the fs, in a yaml format
func NewYamlFileRepo() WorklogRepository {
	return &yamlFileRepo{}
}

type defaults struct {
	duration int `yaml:"duration"`
}

type config struct {
	author   string   `yaml:"author"`
	defaults defaults `yaml:"default"`
}

func (*yamlFileRepo) Configure(author string, duration int) error {
	if err := createDirectory(getWorklogDir()); err != nil {
		return fmt.Errorf("Unable to create directory %s. %s", getWorklogDir(), err.Error())
	}
	file, err := createFile(getWorklogDir() + "config.yml")
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Unable to create configuration file. %s", err.Error())
	}

	config := config{
		author: author,
		defaults: defaults{
			duration: duration,
		},
	}
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("unable to save config. %s", err.Error())
	}
	file.Write(bytes)
	file.Sync()

	return nil
}

func (*yamlFileRepo) Save(wl *model.Work) error {
	fmt.Println("Saving file...")

	file, err := createFile(generateFileName(wl))
	defer file.Close()
	if err != nil {
		return fmt.Errorf("unable to save worklog. %s", err.Error())
	}

	bytes, err := yaml.Marshal(*wl)
	if err != nil {
		return fmt.Errorf("unable to save worklog. %s", err.Error())
	}
	file.Write(bytes)
	file.Sync()

	fmt.Println("Saved file")
	return nil
}

func generateFileName(wl *model.Work) string {
	filePath := getWorklogDir()

	fileName := fmt.Sprintf("%d-%02d-%02d_%02d:%02d_%s",
		wl.When.Year(),
		int(wl.When.Month()),
		wl.When.Day(),
		wl.When.Hour(),
		wl.When.Minute(),
		strings.ReplaceAll(
			strings.TrimSpace(wl.Title), " ", "_"))
	return filePath + fileName + ".yml"
}

func getWorklogDir() string {
	configPath := viper.GetString("default.recordLocation")
	if configPath != "" {
		return configPath
	}

	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return home + "/.worklog/"
}

func createDirectory(filePath string) error {
	err := os.Mkdir(filePath, 0777)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("Unable to create directory '~/.worklog/'. %s", err.Error())
	}
	return nil
}

func createFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to create file %s. %s",
			fileName, err.Error())
	}
	return file, nil
}

func (*yamlFileRepo) GetAllSinceDate(startDate time.Time) ([]*model.Work, error) {
	fmt.Println("Retrieving files...")

	var worklogs []*model.Work
	var errors []string

	fileNames, err := getAllFileNamesFrom(startDate)
	if err != nil {
		return worklogs, err
	}

	for _, fileName := range fileNames {
		readWorklog, err := parseFileToWork(fileName)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			worklogs = append(worklogs, readWorklog)
		}
	}
	if len(errors) != 0 {
		return worklogs, fmt.Errorf("unable to get all files. %s",
			strings.Join(errors, ", "))
	}
	return worklogs, nil
}

func getAllFileNamesFrom(startDate time.Time) ([]string, error) {
	var files []string

	err := filepath.Walk(getWorklogDir(), func(fullPath string, info os.FileInfo, err error) error {
		path := filepath.Base(fullPath)
		if strings.Count(path, "_") < 1 {
			return nil
		}

		splitFileName := strings.Split(path, "_")
		filesDateAsString := fmt.Sprintf("%sT%s:00Z", splitFileName[0], splitFileName[1])
		filesDate, err := helpers.GetStringAsDateTime(filesDateAsString)
		if err != nil {
			return err
		}

		if filesDate.After(startDate) {
			files = append(files, fullPath)
		}
		return nil
	})

	return files, err
}

func parseFileToWork(filePath string) (*model.Work, error) {
	var worklog model.Work
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s. %e", filePath, err)
	}
	err = yaml.Unmarshal(yamlFile, &worklog)
	return &worklog, err
}
