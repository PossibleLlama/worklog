package repository

import (
	"fmt"
	"os"
	"strings"

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
	filePath := viper.GetString("recordLocation")
	if filePath == "" {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		filePath = home + "/.worklog/"
	}

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

func createFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to create file %s. %s",
			fileName, err.Error())
	}
	return file, nil
}
