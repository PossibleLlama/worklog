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
	args := getArguments(os.Args[1:])
	fmt.Printf("%+v\n", args)
}
