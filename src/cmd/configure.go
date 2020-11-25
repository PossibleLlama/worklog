package cmd

import (
	"errors"
	"fmt"

	"github.com/PossibleLlama/worklog/model"

	"github.com/spf13/cobra"
)

var providedAuthor string
var providedDuration int

// configureCmd represents the create command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Setup configuration for the application",
	Long: `Setup configuration file for the application,
setting up the defaults and adding in passed arguments.`,
	Args: func(cmd *cobra.Command, args []string) error {
		providedAuthor = ""
		providedDuration = 0
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return callService()
	},
}

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Default variables to be used",
	Long: `Default variables to be used with the
worklog application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if providedAuthor == "" && providedDuration < 0 {
			return errors.New("defaults requires at least one argument")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return callService()
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.AddCommand(defaultsCmd)

	defaultsCmd.Flags().StringVar(
		&providedAuthor,
		"author",
		"",
		"The authour for all work")
	defaultsCmd.Flags().IntVar(
		&providedDuration,
		"duration",
		-1,
		"Default duration that work takes")
}

func callService() error {
	config := model.NewConfig(providedAuthor, providedDuration)
	if err := wlService.Congfigure(config); err != nil {
		return err
	}
	fmt.Println("Successfully configured")
	return nil
}
