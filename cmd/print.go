package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var startDate time.Time
var startDateString string
var endDate time.Time
var endDateString string

var today bool
var thisWeek bool

var titleFilter string
var descriptionFilter string
var authorFilter string
var rawTagsFilter string
var tagsFilter []string

var prettyOutput bool
var yamlOutput bool
var jsonOutput bool

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print all worklogs since provided date",
	Long: `Prints all worklogs to console that have
been created between the dates provided.`,
	Args: func(cmd *cobra.Command, args []string) error {
		verifySingleFormat()
		verifyFilters()
		return verifyDates()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := model.NewWork(
			titleFilter,
			descriptionFilter,
			authorFilter,
			-1,
			tagsFilter,
			time.Time{})
		worklogs, code, err := wlService.GetWorklogsBetween(startDate, endDate, filter)
		if err != nil {
			return err
		}

		if code == http.StatusNotFound && !jsonOutput {
			fmt.Printf("No work found between %s and %s with the given filter\n",
				startDate, endDate.Add(time.Second*-1))
		} else if prettyOutput {
			model.WriteAllWorkToPrettyText(os.Stdout, worklogs)
		} else if yamlOutput {
			model.WriteAllWorkToPrettyYAML(os.Stdout, worklogs)
		} else {
			model.WriteAllWorkToPrettyJSON(os.Stdout, worklogs)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	// Dates
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

	// Filters
	printCmd.Flags().StringVar(
		&titleFilter,
		"title",
		"",
		"Filter by work including title")
	printCmd.Flags().StringVar(
		&descriptionFilter,
		"description",
		"",
		"Filter by work including description")
	printCmd.Flags().StringVar(
		&authorFilter,
		"author",
		"",
		"Filter by work including author")
	printCmd.Flags().StringVar(
		&rawTagsFilter,
		"tags",
		"",
		"Filter by work including all tags")

	// Format
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

// verifyFilters ensures that the filters make sense
func verifyFilters() {
	titleFilter = strings.TrimSpace(titleFilter)
	descriptionFilter = strings.TrimSpace(descriptionFilter)
	authorFilter = strings.TrimSpace(authorFilter)
	rawTagsList := strings.Split(rawTagsFilter, ",")

	for _, tag := range rawTagsList {
		tagsFilter = append(tagsFilter, strings.TrimSpace(tag))
	}
}

// verifySingleFormat ensures that there is only 1 output format used.
func verifySingleFormat() {
	if !prettyOutput && !yamlOutput && !jsonOutput {
		defaultFormat := viper.GetString("default.format")
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
