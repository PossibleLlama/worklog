package cmd

import (
	"strings"
	"time"

	"github.com/PossibleLlama/worklog/model"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createID string
var createTitle string
var createDescription string
var createWhen time.Time
var createWhenString string
var createDuration int
var createTags []string
var createTagsString string

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
	if createDuration <= -1 {
		createDuration = viper.GetInt("default.duration")
	}
	createTitle = strings.TrimSpace(createTitle)
	createDescription = strings.TrimSpace(createDescription)

	for _, tag := range strings.Split(createTagsString, ",") {
		createTags = append(createTags, strings.TrimSpace(tag))
	}

	whenDate, err := helpers.GetStringAsDateTime(
		strings.TrimSpace(createWhenString))
	if err != nil {
		return err
	}
	createWhen = whenDate

	return nil
}

// CreateRun public method to run create
func CreateRun(cmd *cobra.Command, args []string) error {
	return createRun()
}

func createRun() error {
	wl := model.NewWork(
		createTitle,
		createDescription,
		viper.GetString("author"),
		createDuration,
		createTags,
		createWhen)
	if createID != "" {
		wl.ID = createID
	}
	_, err := wlService.CreateWorklog(wl)
	return err
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVar(
		&createTitle,
		"title",
		"",
		"A short description of the work done")
	createCmd.Flags().StringVar(
		&createDescription,
		"description",
		"",
		"A description of the work")
	createCmd.Flags().StringVar(
		&createWhenString,
		"when",
		helpers.TimeFormat(time.Now()),
		"When the work was worked in RFC3339 format")
	createCmd.Flags().IntVarP(
		&createDuration,
		"duration",
		"",
		-1,
		"Length of time spent on the work")
	createCmd.Flags().StringVar(
		&createTagsString,
		"tags",
		"",
		"Comma seperated list of tags this work relates to")
	createCmd.MarkFlagRequired("title")
}
