package cmd

import (
	"errors"
	"strings"
	"time"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"

	"github.com/spf13/cobra"
)

var (
	editID          string
	editTitle       string
	editDescription string
	editDuration    int
	editAuthor      string
	editWhen        time.Time
	editWhenString  string
	editTags        []string
	editTagsString  string
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing record of work",
	Long:  "Edit an existing record of work by ID, allowing you to update any information previously added. Any provided field will override the existing information.",
	Args:  EditArgs,
	RunE:  EditRun,
}

// EditArgs public method to validate arguments
func EditArgs(cmd *cobra.Command, args []string) error {
	return editArgs(args)
}

func editArgs(args []string) error {
	if len(args) != 1 {
		return errors.New(e.EditID)
	}

	editID = args[0]
	editTitle = strings.TrimSpace(editTitle)
	editDescription = strings.TrimSpace(editDescription)
	editAuthor = strings.TrimSpace(editAuthor)
	editTags = []string{}

	for _, tag := range strings.Split(editTagsString, ",") {
		if strings.TrimSpace(tag) != "" {
			editTags = append(editTags, strings.TrimSpace(tag))
		}
	}

	whenDate, err := helpers.GetStringAsDateTime(editWhenString)
	if err != nil {
		return err
	}
	editWhen = whenDate

	return nil
}

// EditRun public method to run edit
func EditRun(cmd *cobra.Command, args []string) error {
	return editRun(args)
}

func editRun(args []string) error {
	newWl := model.NewWork(
		editTitle,
		editDescription,
		editAuthor,
		editDuration,
		editTags,
		editWhen)
	newWl.ID = editID
	_, err := wlService.EditWorklog(editID, newWl)
	return err
}

func init() {
	rootCmd.AddCommand(editCmd)

	editCmd.Flags().StringVar(
		&editTitle,
		"title",
		"",
		"A short description of the work done")
	editCmd.Flags().StringVar(
		&editDescription,
		"description",
		"",
		"A description of the work")
	editCmd.Flags().StringVar(
		&editAuthor,
		"author",
		"",
		"The author of the work")
	editCmd.Flags().StringVar(
		&editWhenString,
		"when",
		helpers.TimeFormat(time.Now()),
		"When the work was worked in RFC3339 format")
	editCmd.Flags().IntVarP(
		&editDuration,
		"duration",
		"",
		-1,
		"Length of time spent on the work")
	editCmd.Flags().StringVar(
		&editTagsString,
		"tags",
		"",
		"Comma separated list of tags this work relates to")
}
