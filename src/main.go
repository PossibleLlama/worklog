package main

import (
	"fmt"
	"os"
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
