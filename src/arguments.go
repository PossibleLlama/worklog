package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	versionArg      string = "--version"
	versionArgShort string = "-v"
	helpArg         string = "--help"
	helpArgShort    string = "-h"
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
	for _, element := range rawArgs {
		if element == helpArg || element == helpArgShort {
			fmt.Printf(help())
			return emptyArgs, nil
		}
		if element == versionArg || element == versionArgShort {
			fmt.Printf(Version)
			return emptyArgs, nil
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

func help() string {
	return `
		-h --help    | Prints the help function\n
		-v --version | Prints the version of the application\n
	`
}
