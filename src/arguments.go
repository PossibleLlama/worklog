package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	// VersionArg to find the version of the app
	VersionArg string = "--version"
	// HelpArg to find out the options
	HelpArg string = "--help"
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
	emptyArgs := make(map[string]string)

	if len(rawArgs) == 0 {
		return emptyArgs, errors.New("at least one argument is required")
	}
	}
	if len(rawArgs)%2 != 0 {
		return emptyArgs, errors.New("each flag must have an argument")
	}
	for i := 0; i < len(rawArgs); i = i + 2 {
		if rawArgs[i] == "" || rawArgs[i+1] == "" {
			return emptyArgs, errors.New("arguments cannot be empty")
		}
		args[rawArgs[i]] = rawArgs[i+1]
	}

	return args, nil
}
