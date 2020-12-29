package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startDate time.Time
var startDateString string
var endDate time.Time
var endDateString string

var today bool
var thisWeek bool

var yamlOutput bool
var jsonOutput bool
var prettyOutput bool

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print all worklogs since provided date",
	Long: `Prints all worklogs to console that have
been created between the dates provided.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := verifyDates(); err != nil {
			return err
		}
		verifySingleFormat()

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		worklogs, _, err := wlService.GetWorklogsBetween(startDate, endDate)
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
	printCmd.Flags().StringVar(
		&endDateString,
		"endDate",
		"",
		"Date till which to find worklogs. Only functions in conjunction with startDate")
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
	printCmd.Flags().BoolVarP(
		&prettyOutput,
		"pretty",
		"",
		false,
		"Output in a text format")
	printCmd.Flags().BoolVarP(
		&yamlOutput,
		"yaml",
		"",
		false,
		"Output in a yaml format")
	printCmd.Flags().BoolVarP(
		&jsonOutput,
		"json",
		"",
		false,
		"Output in a json format")
}

func verifyDates() error {
	if len(startDateString) != 0 {
		startDateAnytime, err := helpers.GetStringAsDateTime(startDateString)
		if err != nil {
			return err
		}
		startDate = helpers.Midnight(startDateAnytime)
		if len(endDateString) != 0 {
			endDateAnytime, err := helpers.GetStringAsDateTime(endDateString)
			if err != nil {
				return err
			}
			endDate = helpers.Midnight(endDateAnytime).AddDate(0, 0, 1)
		}
	} else if today {
		startDate = helpers.Midnight(time.Now())
		endDate = startDate.AddDate(0, 0, 1)
	} else if thisWeek {
		startDate = helpers.Midnight(helpers.GetPreviousMonday(time.Now()))
		endDate = startDate.AddDate(0, 0, 7)
	} else {
		return errors.New("one flag is required")
	}
	return nil
}

// verifySingleFormat ensures that there is only 1 output format used.
func verifySingleFormat() {
	if !prettyOutput && !yamlOutput && !jsonOutput {
		defaultFormat := viper.GetString("format")
		if defaultFormat == "yaml" || defaultFormat == "yml" {
			yamlOutput = true
		} else if defaultFormat == "json" {
			jsonOutput = true
		} else {
			prettyOutput = true
		}
	} else {
		if prettyOutput {
			yamlOutput = false
			jsonOutput = false
		} else if yamlOutput {
			jsonOutput = false
		}
	}
}
