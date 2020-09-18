package main

import (
	"errors"
	"fmt"
	"os"
)

// Argument is all the possible passed values that will mean something
type Argument string

const (
	// Version to find the version of the app
	Version Argument = "--version"
	// Help argument to find out the options
	Help Argument = "--help"
)

func getArguments(rawArgs []string) map[string]string {
	args, err := formatArguments(rawArgs)
	if err != nil {
		fmt.Printf("passed arguements are not valid. %s\n", err)
		os.Exit(1)
	}
	return args
}

func formatArguments(rawArgs []string) (map[string]string, error) {
	args := make(map[string]string)
	if len(rawArgs) == 0 {
		return args, errors.New("at least one argument is required")
	}
	if len(rawArgs)%2 != 0 {
		return args, errors.New("each argument must have a flag")
	}
	for i := 0; i < len(rawArgs); i = i + 2 {
		if rawArgs[i] == "" || rawArgs[i+1] == "" {
			return args, errors.New("arguments cannot be empty")
		}
	}

	// fmt.Printf("%s: %s\n", rawArgs[0], rawArgs[1])
	// fmt.Printf("%s: %s\n", rawArgs[2], rawArgs[3])

	return args, nil
}
