package cmd

import (
	"errors"
	"fmt"
	"strings"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/model"

	"github.com/spf13/cobra"
)

const (
	configDefaultAuthor   = ""
	configDefaultDuration = 15
	configDefaultFormat   = "pretty"
)

var configProvidedAuthor string
var configProvidedDuration int
var configProvidedFormat string

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
	configProvidedAuthor = configDefaultAuthor
	configProvidedDuration = configDefaultDuration
	configProvidedFormat = configDefaultFormat
	return nil
}

// ConfigRun public method to run configuration
func ConfigRun(cmd *cobra.Command, args []string) error {
	return configRun()
}

func configRun() error {
	cfg := model.NewConfig(configProvidedAuthor, configProvidedFormat, configProvidedDuration)
	if err := wlConfig.SaveConfig(cfg); err != nil {
		return err
	}
	fmt.Println("Successfully configured")
	return nil
}

var overrideDefaultsCmd = &cobra.Command{
	Use:   "overrideDefaults",
	Short: "Override the default variables",
	Long: `Override the default variables to be used with the
worklog application.`,
	Args: OverrideDefaultsArgs,
	RunE: ConfigRun,
}

// OverrideDefaultsArgs public method to validate arguments
func OverrideDefaultsArgs(cmd *cobra.Command, args []string) error {
	return overrideDefaultsArgs()
}

func overrideDefaultsArgs() error {
	configProvidedAuthor = strings.TrimSpace(configProvidedAuthor)
	configProvidedFormat = strings.TrimSpace(configProvidedFormat)
	if configProvidedAuthor == "" &&
		configProvidedFormat == "" &&
		configProvidedDuration < 0 {
		return errors.New(e.ConfigureArgsMinimum)
	}
	if configProvidedDuration < 0 {
		configProvidedDuration = configDefaultDuration
	}
	if configProvidedFormat != "" &&
		configProvidedFormat != "pretty" &&
		configProvidedFormat != "json" &&
		configProvidedFormat != "yaml" {
		return errors.New(e.Format)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.AddCommand(overrideDefaultsCmd)

	overrideDefaultsCmd.Flags().StringVar(
		&configProvidedAuthor,
		"author",
		"",
		"The authour for all work")
	overrideDefaultsCmd.Flags().IntVar(
		&configProvidedDuration,
		"duration",
		-1,
		"Default duration that work takes")
	overrideDefaultsCmd.Flags().StringVar(
		&configProvidedFormat,
		"format",
		"",
		"Format to print work in. If provided, must be one of 'pretty', 'yaml', 'json'")
}
