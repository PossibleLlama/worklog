package cmd

import (
	"fmt"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/spf13/cobra"
)

var title string
var description string
var when time.Time

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new record of work",
	Long: `Creating a new record of work that
the user has created.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("create called with %s %s %s\n", title, description, helpers.TimeFormat(when))
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(
		&title,
		"title",
		"",
		"A short description of the work done")
	createCmd.Flags().StringVar(
		&description,
		"description",
		"",
		"A description of the work")
	createCmd.Flags().String(
		"when",
		helpers.TimeFormat(time.Now()),
		"When the work was worked")
	createCmd.MarkFlagRequired("title")
}
