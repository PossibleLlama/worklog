package main

import (
	"time"
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
		store(work, metadata.RecordLocation)
	}

	print(args, metadata.RecordLocation)
}

func record(args map[string]string, metadata MetadataFile) *Work {
	if title, contained := args[titleArg]; contained {
		fmt.Printf("Saving file...\n")
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

func print(args map[string]string, location string) {
	fmt.Printf("Retrieving file...\n")
	previousDate := getPrintArgumentAsDate(args)
	// TODO search for files with names that are after
	// the YYYY-MM-DD of this date
	fmt.Println(previousDate)
}

func getPrintArgumentAsDate(args map[string]string) time.Time {
	dateShort, containedShort := args[printArgShort]
	dateLong, containedLong := args[printArg]
	var dateString string
	if containedShort || containedLong {
		if containedShort {
			dateString = dateShort
		} else {
			dateString = dateLong
		}
	}
	if len(dateString) == 10 {
		dateString = fmt.Sprintf("%sT00:00:00Z", dateString)
	}
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		fmt.Printf("Date to print from is not a valid date. %s\n", err.Error())
		os.Exit(1)
	}
	return date
}
