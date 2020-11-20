package repository

import (
	"fmt"

	"github.com/PossibleLlama/worklog/model"
)

type yamlFileRepo struct{}

// NewYamlFileRepo Generator for repository storing worklogs
// on the fs, in a yaml format
func NewYamlFileRepo() WorklogRepository {
	return &yamlFileRepo{}
}

func (*yamlFileRepo) Save(wl *model.Work) error {
	fmt.Println("Saving file...")
	// Do stuff
	fmt.Println("Saved file")
	return nil
}
