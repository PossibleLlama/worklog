package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/repository"
	"github.com/PossibleLlama/worklog/service"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	wlService service.WorklogService
	wlRepo    repository.WorklogRepository
)
var (
	home         string
	cfgFile      string
	repoLocation string
)

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

// Execute adds all child commands to the root command and sets flags appropriately
// This is called by main.main(). It only needs to happen once to the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(e.StartupErrors)
	}
}

func init() {
	var err error

	home, err = homedir.Dir()
	if err != nil {
		fmt.Println("Unable to get home directory", err)
		os.Exit(e.StartupErrors)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		fmt.Sprintf("%s%s.worklog%sconfig.yml",
			home,
			string(filepath.Separator),
			string(filepath.Separator)),
		"config file including file extension (default is $HOME/.worklog/config.yml)")
	rootCmd.PersistentFlags().StringVar(&repoLocation,
		"repo",
		fmt.Sprintf("%s%s.worklog%sworklog.db",
			home,
			string(filepath.Separator),
			string(filepath.Separator)),
		"repository that worklogs are stored in (default is $HOME/.worklog/worklog.db)")

	cobra.OnInitialize(initConfig)

	wlRepo = repository.NewYamlFileRepo()
	wlService = service.NewWorklogService(wlRepo)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if !strings.HasPrefix(cfgFile, string(filepath.Separator)) {
		cfgFile = fmt.Sprintf("%s%s%s", home, string(filepath.Separator), cfgFile)
	}

	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Unable to use config file: '%s'. %s\n", viper.ConfigFileUsed(), err)
		os.Exit(e.StartupErrors)
	}
}
