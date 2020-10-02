package main

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"os"
	"strings"
)

const (
	// Version keep track of the version of the application
	Version string = "worklog-0.0.1"
)

func main() {
	metadataChan := make(chan MetadataFile)

	go getMetadata(metadataChan)
	args := getArguments(os.Args[1:])

	metadata := <-metadataChan

	work := record(args, metadata)
	if work != nil {
		fmt.Printf("%s\n", *work)
		store(work, metadata.RecordLocation)
	}
}

func record(args map[string]string, metadata MetadataFile) *Work {
	if title, contained := args[titleArg]; contained {
		return New(title,
			args[descriptionArg],
			metadata.Author,
			args[whereArg],
			args[whenArg])
	}
	return nil
}

func store(work *Work, location string) {
	fileName := fmt.Sprintf("%d-%02d-%02d_%02d:%02d_%s",
		work.When.Year(),
		int(work.When.Month()),
		work.When.Day(),
		work.When.Hour(),
		work.When.Minute(),
		strings.ReplaceAll(work.Title, " ", "_"))
	
	file, err := os.Create(location + fileName)
	if err != nil {
		fmt.Printf("Unable to create file %s. %s\n",
			location + work.When.String(),
			err.Error())
		os.Exit(1)
	}
	defer file.Close()

	bytes, err := yaml.Marshal(*work)
	if err != nil {
		fmt.Printf("Unable to encode record. %s\n", work)
		os.Exit(1)
	}
	file.Write(bytes)
	file.Sync()
}
