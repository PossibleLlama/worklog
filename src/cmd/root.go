package cmd

import (
	"fmt"
	"os"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/repository"

	"github.com/spf13/cobra"

	"github.com/PossibleLlama/worklog/service"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var wlService service.WorklogService
var wlRepo repository.WorklogRepository
var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "worklog",
	Version: helpers.Version,
	Short:   "A productivity tool to track previous work",
	Long: `A CLI tool to let people track what work they
have completed. You can record what work you do,
and get a summary of what you've done each day.

For information on using the CLI, use worklog
--help`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	wlRepo = repository.NewYamlFileRepo()
	wlService = service.NewWorklogService(wlRepo)

	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		"",
		"config file (default is $HOME/.worklog/config.yml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".worklog" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".worklog/config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Unable to use config file:", viper.ConfigFileUsed())
	}
}
