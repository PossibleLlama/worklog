package main

import "github.com/PossibleLlama/worklog/cmd"

const (
	// Version keep track of the version of the application
	Version string = "worklog-0.1.0"
)

func main() {
	cmd.Execute()
}
