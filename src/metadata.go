package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	configFileName	string = "config.yml"
)

// MetadataFile information added to each worklog
type MetadataFile struct {
	Author			string `yaml:"author"`
	RecordLocation	string `yaml:"recordLocation"`
}

func getMetadata(metadataChan chan<- MetadataFile) {
	var file MetadataFile
	filePath := os.Getenv("HOME") + "/.worklog/"
	configFilePath := filePath + configFileName

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Printf("cannot open file: %s\n", configFilePath)
		metadataChan <- MetadataFile{}
	}

	err = yaml.Unmarshal(data, &file)
	if err != nil {
		fmt.Printf("cannot unmarshal data: %s\n", err)
		metadataChan <- MetadataFile{}
	}

	if file.Author == "" {
		fmt.Printf("unable to get 'author' property from ~%s file.\n", configFileName)
		os.Exit(1)
	}
	if file.RecordLocation == "" {
		file.RecordLocation = filePath
	}

	metadataChan <- file
}
