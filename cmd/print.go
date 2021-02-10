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

var printStartDate time.Time
var printStartDateString string
var printEndDate time.Time
var printEndDateString string

var printToday bool
var printThisWeek bool

var printFilterTitle string
var printFilterDescription string
var printFilterAuthor string
var printFilterTags []string
var printFilterTagsString string

var printOutputPretty bool
var printOutputYAML bool
var printOutputJSON bool

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print all worklogs since provided date",
	Long: `Prints all worklogs to console that have
been created between the dates provided.`,
	Args: PrintArgs,
	RunE: PrintRun,
}

// PrintArgs public method to validate arguments
func PrintArgs(cmd *cobra.Command, args []string) error {
	return printArgs(args)
}

func printArgs(args []string) error {
	verifySingleFormat()
	verifyFilters()
	return verifyDatesAndIDs(args)
}

// PrintRun public method to run print
func PrintRun(cmd *cobra.Command, args []string) error {
	return printRun(args...)
}

func printRun(args ...string) error {
	// Passing args through to allow for specifying ID's
	filter := &model.Work{
		Title:       printFilterTitle,
		Description: printFilterDescription,
		Author:      printFilterAuthor,
		Duration:    -1,
		Tags:        printFilterTags,
		When:        time.Time{},
		CreatedAt:   time.Time{}}
	worklogs, code, err := wlService.GetWorklogsBetween(printStartDate, printEndDate, filter)
	if err != nil {
		return err
	}

	if code == http.StatusNotFound && !printOutputJSON {
		fmt.Printf("No work found between %s and %s with the given filter\n",
			printStartDate, printEndDate.Add(time.Second*-1))
	} else if printOutputPretty {
		model.WriteAllWorkToPrettyText(os.Stdout, worklogs)
	} else if printOutputYAML {
		model.WriteAllWorkToPrettyYAML(os.Stdout, worklogs)
	} else {
		model.WriteAllWorkToPrettyJSON(os.Stdout, worklogs)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(printCmd)

	// Dates
	printCmd.Flags().StringVar(
		&printStartDateString,
		"startDate",
		"",
		"Date from which to find worklogs")
	printCmd.Flags().StringVar(
		&printEndDateString,
		"endDate",
		"",
		"Date till which to find worklogs. Only functions in conjunction with startDate")
	printCmd.Flags().BoolVarP(
		&printToday,
		"today",
		"",
		false,
		"Print today's work")
	printCmd.Flags().BoolVarP(
		&printThisWeek,
		"thisWeek",
		"",
		false,
		"Prints this weeks work")

	// Filters
	printCmd.Flags().StringVar(
		&printFilterTitle,
		"title",
		"",
		"Filter by work including title")
	printCmd.Flags().StringVar(
		&printFilterDescription,
		"description",
		"",
		"Filter by work including description")
	printCmd.Flags().StringVar(
		&printFilterAuthor,
		"author",
		"",
		"Filter by work including author")
	printCmd.Flags().StringVar(
		&printFilterTagsString,
		"tags",
		"",
		"Filter by work including all tags")

	// Format
	printCmd.Flags().BoolVarP(
		&printOutputPretty,
		"pretty",
		"",
		false,
		"Output in a text format")
	printCmd.Flags().BoolVarP(
		&printOutputYAML,
		"yaml",
		"",
		false,
		"Output in a yaml format")
	printCmd.Flags().BoolVarP(
		&printOutputJSON,
		"json",
		"",
		false,
		"Output in a json format")
}

// verifySingleFormat ensures that there is only 1 output format used.
func verifySingleFormat() {
	if !printOutputPretty && !printOutputYAML && !printOutputJSON {
		defaultFormat := viper.GetString("default.format")
		if defaultFormat == "yaml" || defaultFormat == "yml" {
			printOutputYAML = true
		} else if defaultFormat == "json" {
			printOutputJSON = true
		} else {
			printOutputPretty = true
		}
	} else {
		if printOutputPretty {
			printOutputYAML = false
			printOutputJSON = false
		} else if printOutputYAML {
			printOutputJSON = false
		}
	}
}

// verifyFilters ensures that the filters make sense
func verifyFilters() {
	printFilterTitle = strings.TrimSpace(printFilterTitle)
	printFilterDescription = strings.TrimSpace(printFilterDescription)
	printFilterAuthor = strings.TrimSpace(printFilterAuthor)
	rawTagsList := strings.Split(printFilterTagsString, ",")

	for _, tag := range rawTagsList {
		if strings.TrimSpace(tag) != "" {
			printFilterTags = append(printFilterTags, strings.TrimSpace(tag))
		}
	}
}

func verifyDatesAndIDs(ids []string) error {
	errID := verifyIDs(ids)
	errDates := verifyDates()
	if errDates == nil {
		return nil
	} else if errID == nil {
		return nil
	}
	return errDates
}

func verifyIDs(ids []string) error {
	if len(ids) == 0 {
		return errors.New("No ids provided")
	}
	return nil
}

// verifyDates ensures the dates are valid
func verifyDates() error {
	if len(printStartDateString) != 0 {
		startDateAnytime, err := helpers.GetStringAsDateTime(printStartDateString)
		if err != nil {
			return err
		}
		printStartDate = helpers.Midnight(startDateAnytime)
		if len(printEndDateString) != 0 {
			endDateAnytime, err := helpers.GetStringAsDateTime(printEndDateString)
			if err != nil {
				return err
			}
			printEndDate = helpers.Midnight(endDateAnytime).AddDate(0, 0, 1)
		}
	} else if printToday {
		printStartDate = helpers.Midnight(time.Now())
		printEndDate = printStartDate.AddDate(0, 0, 1)
	} else if printThisWeek {
		printStartDate = helpers.GetPreviousMonday(time.Now())
		printEndDate = printStartDate.AddDate(0, 0, 7)
	} else {
		return errors.New("one flag is required")
	}
	return nil
}
