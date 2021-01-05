package cmd

import (
	"strings"
	"time"

	"github.com/PossibleLlama/worklog/model"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var title string
var description string
var when time.Time
var whenString string
var duration int
var tags []string
var tagsString string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new record of work",
	Long: `Creating a new record of work that
the user has created.`,
	Args: CreateArgs,
	RunE: CreateRun,
}

// CreateArgs public method to validate arguments
func CreateArgs(cmd *cobra.Command, args []string) error {
	return createArgs()
}

func createArgs() error {
	if duration <= -1 {
		duration = viper.GetInt("default.duration")
	}
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	whenDate, err := helpers.GetStringAsDateTime(
		strings.TrimSpace(whenString))
	if err != nil {
		return err
	}
	when = whenDate

	for _, tag := range strings.Split(tagsString, ",") {
		tags = append(tags, strings.TrimSpace(tag))
	}

	return nil
}

// CreateRun public method to run create
func CreateRun(cmd *cobra.Command, args []string) error {
	return createRun()
}

func createRun() error {
	_, err := wlService.CreateWorklog(model.NewWork(
		title,
		description,
		viper.GetString("author"),
		duration,
		tags,
		when))
	return err
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
	createCmd.Flags().IntVarP(
		&duration,
		"duration",
		"",
		-1,
		"Length of time spent on the work")
	createCmd.Flags().StringVar(
		&tagsString,
		"tags",
		"",
		"Comma seperated list of tags this work relates to")
	createCmd.MarkFlagRequired("title")
}
