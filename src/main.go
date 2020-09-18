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

	fmt.Printf("%+v\n%+v\n", args, metadata)
}
