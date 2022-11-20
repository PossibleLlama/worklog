package cli

import (
	"os"
	"strings"
	"time"

	"github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	createID          string
	createTitle       string
	createDescription string
	createAuthor      string
	createWhen        time.Time
	createWhenString  string
	createDuration    int
	createTags        []string
	createTagsString  string
)

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
	createTitle = helpers.Sanitize(strings.TrimSpace(createTitle))
	createDescription = helpers.Sanitize(strings.TrimSpace(createDescription))
	createAuthor = helpers.Sanitize(strings.TrimSpace(createAuthor))

	if createAuthor == "" {
		createAuthor = helpers.Sanitize(viper.GetString("default.author"))
	}

	for _, tag := range strings.Split(createTagsString, ",") {
		if strings.TrimSpace(tag) != "" {
			createTags = append(createTags, helpers.Sanitize(strings.TrimSpace(tag)))
		}
	}

	whenDate, err := helpers.GetStringAsDateTime(createWhenString)
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
		createAuthor,
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
		&createAuthor,
		"author",
		"",
		"The author of the work")
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
		"Comma separated list of tags this work relates to")
	if err := createCmd.MarkFlagRequired("title"); err != nil {
		os.Exit(errors.StartupErrors)
	}
}
