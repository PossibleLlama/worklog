package cmd

import (
	"time"

	"github.com/PossibleLlama/worklog/model"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/spf13/cobra"
)

var title string
var description string
var when time.Time
var whenString string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new record of work",
	Long: `Creating a new record of work that
the user has created.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		when, err := helpers.GetStringAsDateTime(whenString)
		if err != nil {
			return err
		}

		work := model.NewWork(title, description, "", "", helpers.TimeFormat(when))
		return wlService.CreateWorklog(work)
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
	createCmd.Flags().StringVar(
		&whenString,
		"when",
		helpers.TimeFormat(time.Now()),
		"When the work was worked in RFC3339 format")
	createCmd.MarkFlagRequired("title")
}
