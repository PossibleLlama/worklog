package cli

import (
	"errors"
	"strings"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"

	"github.com/spf13/cobra"
)

const (
	configDefaultAuthor   = ""
	configDefaultDuration = 15
	configDefaultFormat   = "pretty"
	configDefaultRepoType = "bolt"
	configDefaultRepoPath = ""
)

var (
	configProvidedAuthor   string
	configProvidedDuration int
	configProvidedFormat   string
	configProvidedRepoType string
	configProvidedRepoPath string
)

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
	configProvidedRepoType = configDefaultRepoType
	configProvidedRepoPath = configDefaultRepoPath
	return nil
}

// ConfigRun public method to run configuration
func ConfigRun(cmd *cobra.Command, args []string) error {
	return configRun()
}

func configRun() error {
	cfg := model.NewConfig(
		model.Defaults{
			Author:   configProvidedAuthor,
			Format:   configProvidedFormat,
			Duration: configProvidedDuration,
		}, model.Repo{
			Type: configProvidedRepoType,
			Path: configProvidedRepoPath,
		})
	if err := wlConfig.SaveConfig(cfg); err != nil {
		return err
	}
	helpers.LogInfo("Successfully configured", "configure - saved config")
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
	configProvidedRepoType = strings.TrimSpace(configProvidedRepoType)
	configProvidedRepoPath = strings.TrimSpace(configProvidedRepoPath)
	if configProvidedAuthor == "" &&
		configProvidedFormat == "" &&
		configProvidedDuration < 0 &&
		configProvidedRepoType == "" &&
		configProvidedRepoPath == "" {
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
	if configProvidedRepoType != "" &&
		configProvidedRepoType != "bolt" &&
		configProvidedRepoType != "legacy" {
		return errors.New(e.RootRepoType)
	}
	// Repo path will accept anything, it's up to the user to make sure
	// the file path makes sense.
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
	overrideDefaultsCmd.Flags().StringVar(
		&configProvidedRepoType,
		"repo",
		"bolt",
		"The type of repository used")
	overrideDefaultsCmd.Flags().StringVar(
		&configProvidedRepoPath,
		"repoPath",
		"",
		"The path to the repository for storage")
}
