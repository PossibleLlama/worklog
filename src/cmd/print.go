package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/spf13/cobra"
)

var startDate time.Time
var startDateString string
var today bool
var thisWeek bool

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print all worklogs since provided date",
	Long: `Prints all worklogs to console that have
been created since the start provided date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(startDateString) != 0 {
			startDateAnytime, err := helpers.GetStringAsDateTime(startDateString)
			if err != nil {
				return err
			}
			startDate = helpers.Midnight(startDateAnytime)
		} else if today {
			startDate = helpers.Midnight(time.Now())
		} else if thisWeek {
			startDate = helpers.Midnight(helpers.GetPreviousMonday(time.Now()))
		} else {
			return errors.New("require one flag to be provided")
		}

		worklogs, _, err := wlService.GetWorklogsSince(startDate)
		if err != nil {
			return err
		}

		for _, work := range worklogs {
			fmt.Printf("%+v\n", work)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	printCmd.Flags().StringVar(
		&startDateString,
		"startDate",
		"",
		"Date from which to find worklogs")
	printCmd.Flags().BoolVarP(
		&today,
		"today",
		"",
		false,
		"Print today's work")
	printCmd.Flags().BoolVarP(
		&thisWeek,
		"thisWeek",
		"",
		false,
		"Prints this weeks work")
}
