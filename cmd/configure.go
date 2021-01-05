package cmd

import (
	"errors"
	"fmt"

	"github.com/PossibleLlama/worklog/model"

	"github.com/spf13/cobra"
)

const defaultAuthor = ""
const defaultDuration = 15

var providedAuthor string
var providedDuration int
var providedFormat string

// configureCmd represents the create command
var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Setup configuration for the application",
	Long: `Setup configuration file for the application,
setting up the defaults and adding in passed arguments.`,
	Args: ConfigArgs,
	RunE: ConfigRun,
}

// ConfigArgs public method to validate arguments
func ConfigArgs(cmd *cobra.Command, args []string) error {
	return configArgs()
}

func configArgs() error {
	providedAuthor = defaultAuthor
	providedDuration = defaultDuration
	return nil
}

// ConfigRun public method to run configuration
func ConfigRun(cmd *cobra.Command, args []string) error {
	return configRun()
}

func configRun() error {
	cfg := model.NewConfig(providedAuthor, providedFormat, providedDuration)
	if err := wlService.Congfigure(cfg); err != nil {
		return err
	}
	fmt.Println("Successfully configured")
	return nil
}

var defaultsCmd = &cobra.Command{
	Use:   "defaults",
	Short: "Default variables to be used",
	Long: `Default variables to be used with the
worklog application.`,
	Args: DefaultArgs,
	RunE: ConfigRun,
}

// DefaultArgs public method to validate arguments
func DefaultArgs(cmd *cobra.Command, args []string) error {
	return defaultArgs(args...)
}

func defaultArgs(args ...string) error {
	if providedAuthor == "" &&
		providedFormat == "" &&
		providedDuration < 0 {
		return errors.New("defaults requires at least one argument")
	}
	if providedDuration < 0 {
		providedDuration = defaultDuration
	}
	if providedFormat != "" &&
		providedFormat != "pretty" &&
		providedFormat != "json" &&
		providedFormat != "yaml" &&
		providedFormat != "yml" {
		return errors.New("provided format is not valid")
	}
	return nil
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
	defaultsCmd.Flags().StringVar(
		&providedFormat,
		"format",
		"",
		"Format to print work in. If provided, must be one of 'pretty', 'yaml', 'json'")
}
