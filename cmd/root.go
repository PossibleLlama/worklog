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
	wlConfig  repository.ConfigRepository
)
var (
	homeDir      string
	cfgFile      string
	repoType     string
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

	homeDir, err = homedir.Dir()
	if err != nil {
		fmt.Println("Unable to get home directory", err)
		os.Exit(e.StartupErrors)
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		// This does not contain the home directory, as then the tool tip description
		// would include that value, which is difficult to test for
		fmt.Sprintf(".worklog%sconfig.yml",
			string(filepath.Separator)),
		"Config file including file extension")
	rootCmd.PersistentFlags().StringVar(&repoType,
		"repo",
		"legacy",
		"Which type of repository to use for storing/retrieving worklogs")
	rootCmd.PersistentFlags().StringVar(&repoLocation,
		"repoPath",
		// This does not contain the home directory, as then the tool tip description
		// would include that value, which is difficult to test for
		fmt.Sprintf(".worklog%sworklog.db",
			string(filepath.Separator)),
		"Directory path for repository that worklogs are stored in")

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if !strings.HasPrefix(cfgFile, string(filepath.Separator)) {
		cfgFile = fmt.Sprintf("%s%s%s", homeDir, string(filepath.Separator), cfgFile)
	}
	if !strings.HasPrefix(repoLocation, string(filepath.Separator)) {
		repoLocation = fmt.Sprintf("%s%s%s", homeDir, string(filepath.Separator), repoLocation)
	}

	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Unable to use config file: '%s'. %s\n", viper.ConfigFileUsed(), err)
	}

	if strings.ToLower(repoType) == "legacy" {
		wlRepo = repository.NewYamlFileRepo()
	} else if strings.ToLower(repoType) == "bolt" {
		wlRepo = repository.NewBBoltRepo(repoLocation)
	} else {
		fmt.Printf("A valid repo must be specified. 'legacy' or 'bolt'\n")
		os.Exit(e.StartupErrors)
	}
	wlConfig = repository.NewYamlConfig(
		filepath.Dir(cfgFile))
	wlService = service.NewWorklogService(wlRepo)
}
