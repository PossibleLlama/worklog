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
var tags string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new record of work",
	Long: `Creating a new record of work that
the user has created.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if duration <= -1 {
			duration = viper.GetInt("default.duration")
		}
		title = strings.TrimSpace(title)
		description = strings.TrimSpace(description)
		whenString = strings.TrimSpace(whenString)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		when, err := helpers.GetStringAsDateTime(whenString)
		if err != nil {
			return err
		}

		var tagsList []string
		for _, tag := range strings.Split(tags, ",") {
			tagsList = append(tagsList, strings.TrimSpace(tag))
		}

		work := model.NewWork(
			title,
			description,
			viper.GetString("author"),
			duration,
			tagsList,
			when)
		_, err = wlService.CreateWorklog(work)
		return err
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
	createCmd.Flags().IntVarP(
		&duration,
		"duration",
		"",
		-1,
		"Length of time spent on the work")
	createCmd.Flags().StringVar(
		&tags,
		"tags",
		"",
		"Comma seperated list of tags this work relates to")
	createCmd.MarkFlagRequired("title")
}
