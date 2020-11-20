package cmd

import (
	"fmt"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/spf13/cobra"
)

var startDate time.Time
var startDateString string

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print all worklogs since provided date",
	Long: `Prints all worklogs to console that have
been created since the start provided date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		startDateAnytime, err := helpers.GetStringAsDateTime(startDateString)
		startDate = helpers.Midnight(startDateAnytime)
		if err != nil {
			return err
		}

		worklogs, _, err := wlService.GetWorklogsSince(startDate)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", worklogs)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	printCmd.Flags().StringVar(
		&startDateString,
		"startDate",
		helpers.TimeFormat(time.Now()),
		"Date from which to find worklogs")
	printCmd.MarkFlagRequired("startDate")
}
