package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	fileName string = "/.worklog.yml"
)

// MetadataFile information added to each worklog
type MetadataFile struct {
	Author string `yaml:"author"`
}

func getMetadata(metadataChan chan<- MetadataFile) {
	var file MetadataFile
	filePath := os.Getenv("HOME") + fileName

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("cannot open file: %s\n", filePath)
		metadataChan <- MetadataFile{}
	}

	err = yaml.Unmarshal(data, &file)
	if err != nil {
		fmt.Printf("cannot unmarshal data: %s\n", err)
		metadataChan <- MetadataFile{}
	}

	if file.Author == "" {
		fmt.Printf("unable to get 'author' property from ~%s file.\n", fileName)
		os.Exit(1)
	}

	metadataChan <- file
}
