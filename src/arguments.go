package main

import (
	"errors"
	"fmt"
	"os"
)

const (
	helpArg         string = "--help"
	helpArgShort    string = "-h"
	versionArg      string = "--version"
	versionArgShort string = "-v"
	printArg        string = "--print"
	printArgShort   string = "-p"
	recordArg       string = "--record"
	recordArgShort  string = "-r"
)

func listAllSingleArgs() []string {
	return []string{
		helpArg,
		helpArgShort,
		versionArg,
		versionArgShort,
	}
}

func listAllPairArgs() []string {
	return []string{
		printArg,
		printArgShort,
		recordArg,
		recordArgShort,
	}
}

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
	if isValidSingleArgument(rawArgs) {
		return emptyArgs, nil
	} else if len(rawArgs)%2 != 0 {
		return emptyArgs, errors.New("each flag must have an argument")
	}

	for i := 0; i < len(rawArgs); i = i + 2 {
		if rawArgs[i] == "" || rawArgs[i+1] == "" {
			return emptyArgs, errors.New("arguments cannot be empty")
		}

		if !isValidPairArgument(rawArgs[i]) {
			fmt.Printf(help())
			return emptyArgs, nil
		}
		args[rawArgs[i]] = rawArgs[i+1]
	}

	return args, nil
}

func isValidSingleArgument(rawArgs []string) bool {
	for _, element := range rawArgs {
		if element == helpArg || element == helpArgShort {
			fmt.Printf(help())
			os.Exit(0)
		}
		if element == versionArg || element == versionArgShort {
			fmt.Printf("%s\n", Version)
			os.Exit(0)
		}
	}
	return false
}

func isValidPairArgument(arg string) bool {
	for _, element := range listAllPairArgs() {
		if arg == element {
			return true
		}
	}
	return false
}

func help() string {
	return " Argument     | Description\n" +
		"----------------------------------------\n" +
		"-h --help     | Prints the help function.\n" +
		"-v --version  | Prints the version of the application.\n" +
		"-r --record   | Makes a record of work done.\n" +
		"-p --print    | Prints the work done.\n"
}
