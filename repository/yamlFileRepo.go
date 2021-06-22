package repository

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"

	homedir "github.com/mitchellh/go-homedir"
)

type yamlFileRepo struct{}

// NewYamlFileRepo Generator for repository storing worklogs
// on the fs, in a yaml format
func NewYamlFileRepo() WorklogRepository {
	return &yamlFileRepo{}
}

func (*yamlFileRepo) Configure(cfg *model.Config) error {
	if err := createDirectory(getWorklogDir()); err != nil {
		return fmt.Errorf("%s %s. %s", e.RepoCreateDirectory, getWorklogDir(), err.Error())
	}
	file, err := createFile(getWorklogDir() + "config.yml")
	if err != nil {
		return fmt.Errorf("%s. %s", e.RepoConfigFileCreate, err.Error())
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	if err := cfg.WriteYAML(file); err != nil {
		return fmt.Errorf("%s. %s", e.RepoConfigFileSave, err.Error())
	}
	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}

func (*yamlFileRepo) Save(wl *model.Work) error {
	fmt.Println("Saving file...")

	file, err := createFile(generateFileName(wl))
	if err != nil {
		return fmt.Errorf("%s. %s", e.RepoSaveFile, err.Error())
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}()

	if err := wl.WriteYAML(file); err != nil {
		return fmt.Errorf("%s. %s", e.RepoSaveFile, err.Error())
	}
	if err := file.Sync(); err != nil {
		return err
	}

	fmt.Println("Saved file")
	return nil
}

func generateFileName(wl *model.Work) string {
	filePath := getWorklogDir()

	fileName := fmt.Sprintf("%d-%02d-%02dT%02d:%02d_%d_%s",
		wl.When.Year(),
		int(wl.When.Month()),
		wl.When.Day(),
		wl.When.Hour(),
		wl.When.Minute(),
		wl.Revision,
		wl.ID)
	return filePath + fileName + ".yml"
}

func getWorklogDir() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(e.RepoErrors)
	}
	return home + "/.worklog/"
}

func createDirectory(filePath string) error {
	err := os.Mkdir(filePath, 0750)
	if err != nil {
		if os.IsExist(err) {
			return nil
		}
		return fmt.Errorf("%s '%s'. %s", e.RepoCreateDirectory, filePath, err.Error())
	}
	return nil
}

func createFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s %s. %s", e.RepoCreateFile,
			fileName, err.Error())
	}
	return file, nil
}

func (*yamlFileRepo) GetAllBetweenDates(startDate, endDate time.Time, filter *model.Work) ([]*model.Work, error) {
	var worklogs []*model.Work
	var errors []string

	if (endDate == time.Time{}) {
		endDate = time.Date(3000, time.January, 1, 0, 0, 0, 0, time.Now().Location())
	}

	fileNames, err := getAllFileNamesBetweenDates(startDate, endDate)
	if err != nil {
		return worklogs, err
	}

	for _, fileName := range fileNames {
		readWorklog, err := parseFileToWork(fileName)
		if err != nil {
			errors = append(errors, err.Error())
		} else if workMatchesFilter(filter, readWorklog) {
			worklogs = append(worklogs, readWorklog)
		}
	}
	if len(errors) != 0 {
		return worklogs, fmt.Errorf("%s. %s", e.RepoGetFiles,
			strings.Join(errors, ", "))
	}
	return worklogs, nil
}

func (*yamlFileRepo) GetByID(ID string, filter *model.Work) (*model.Work, error) {
	var wl *model.Work
	var err error

	fileName, err := getFileByID(ID)
	if err != nil {
		return nil, err
	}

	if fileName == "" {
		return nil, nil
	}
	wl, err = parseFileToWork(fileName)
	if err != nil {
		return nil, err
	} else if workMatchesFilter(filter, wl) {
		return wl, nil
	}
	return nil, nil
}

func getAllFileNamesBetweenDates(startDate, endDate time.Time) ([]string, error) {
	var files []string

	err := filepath.Walk(getWorklogDir(), func(fullPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		path := filepath.Base(fullPath)
		if strings.Count(path, "_") < 2 {
			return nil
		}

		splitFileName := strings.Split(path, "_")
		filesDateAsString := fmt.Sprintf("%s:00Z", splitFileName[0])
		filesDate, err := helpers.GetStringAsDateTime(filesDateAsString)
		if err != nil {
			return err
		}

		if filesDate.After(startDate) && filesDate.Before(endDate) {
			files = append(files, fullPath)
		}
		return nil
	})

	return files, err
}

func getFileByID(ID string) (string, error) {
	ids := make(map[string]string)

	err := filepath.Walk(getWorklogDir(), func(fullPath string, info os.FileInfo, err error) error {
		path := filepath.Base(fullPath)
		if strings.Count(path, "_") < 2 {
			return nil
		}

		splitFileName := strings.Split(path, "_")
		if aInB(ID, splitFileName[2]) {
			currentRev, err := strconv.Atoi(splitFileName[1])
			if err != nil {
				return err
			} else if highestRevPath, ok := ids[splitFileName[2]]; ok {
				highestRev, err := strconv.Atoi(strings.Split(highestRevPath, "_")[1])
				if err != nil {
					return err
				} else if currentRev > highestRev {
					ids[splitFileName[2]] = fullPath
				}
			} else {
				ids[splitFileName[2]] = fullPath
			}
		}
		return nil
	})

	if err != nil || len(ids) == 0 {
		return "", err
	} else if len(ids) > 1 {
		return "", fmt.Errorf("ID '%s' is not unique", ID)
	}

	for _, v := range ids {
		return v, nil
	}
	return "", errors.New(e.Unexpected)
}

func parseFileToWork(filePath string) (*model.Work, error) {
	yamlFile, err := ioutil.ReadFile(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("%s %s. %e", e.RepoGetFilesRead, filePath, err)
	}
	worklog, err := model.ReadYAML(yamlFile)
	if err == nil {
		sort.Strings(worklog.Tags)
	}
	return worklog, err
}

func workMatchesFilter(filter, w *model.Work) bool {
	if !aInB(filter.Title, w.Title) {
		return false
	}
	if !aInB(filter.Description, w.Description) {
		return false
	}
	if !aInB(filter.Author, w.Author) {
		return false
	}
	for _, filtersTag := range filter.Tags {
		if filtersTag != "" &&
			!aInB(filtersTag, strings.Join(w.Tags, " ")) {
			return false
		}
	}
	return true
}

func aInB(a, b string) bool {
	return a == "" || strings.Contains(
		strings.ToLower(b),
		strings.ToLower(a))
}
