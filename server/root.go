package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	e "github.com/PossibleLlama/worklog/errors"
	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/repository"
	"github.com/PossibleLlama/worklog/service"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	PATH    = "/worklog"
	ID_PATH = PATH + "/{id}"
)

var (
	httpRouter *mux.Router = mux.NewRouter().StrictSlash(true)
	wlService  service.WorklogService
	wlRepo     repository.WorklogRepository
	wlConfig   repository.ConfigRepository
)

var (
	homeDir      string
	cfgFile      string
	repoType     string
	repoLocation string
	port         int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "worklog-server",
	Version: helpers.Version,
	Short:   "Launching a server for a productivity tool to track work",
	Long: `A server to let people track what work they
have completed. You can record what work you do,
and get a summary of what you've done each day.

For information on using the CLI, use worklog-server
--help`,
}

func Execute() {
	InitCobra()
	InitConfig()
	if err := rootCmd.Execute(); err != nil {
		helpers.LogError(err.Error(), "startup")
		os.Exit(e.StartupErrors)
	}

	startServer()
}

func startServer() {
	helpers.LogInfo(fmt.Sprintf("server starting on port %d\n", port), "startup")

	httpRouter.Use(SetDefaultHeadersMiddleware)

	httpRouter.NotFoundHandler = http.HandlerFunc(NotFound)
	httpRouter.MethodNotAllowedHandler = http.HandlerFunc(InvalidMethod)
	httpRouter.HandleFunc(PATH, Create).Methods(http.MethodPost)
	httpRouter.HandleFunc(PATH, Print).Methods(http.MethodGet)
	httpRouter.HandleFunc(ID_PATH, PrintSingle).Methods(http.MethodGet)
	httpRouter.HandleFunc(ID_PATH, Edit).Methods(http.MethodPut)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut},
	})

	server := &http.Server{
		Addr:              fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout:      time.Second * 10,
		ReadTimeout:       time.Second * 5,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 20,
		Handler:           c.Handler(httpRouter),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			helpers.LogError(fmt.Sprintf("server error occurred '%s'\n", err.Error()), "startup")
			os.Exit(e.StartupErrors)
		}
	}()

	interruptAndExit(server)
}

func interruptAndExit(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)
	<-interruptChan

	helpers.LogInfo("server stopping...", "shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// #nosec G104 -- There could be an error while shutting down, but its during shutdown
	server.Shutdown(ctx)
}

func InitCobra() {
	var err error

	homeDir, err = os.UserHomeDir()
	if err != nil {
		helpers.LogError(fmt.Sprintf("Unable to get home directory: %s", err.Error()), "startup")
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
		"",
		"Which type of repository to use for storing/retrieving worklogs")
	rootCmd.PersistentFlags().StringVar(&repoLocation,
		"repoPath",
		// This does not contain the home directory, as then the tool tip description
		// would include that value, which is difficult to test for
		"",
		"Directory path for repository that worklogs are stored in")
	rootCmd.PersistentFlags().IntVarP(&port,
		"port",
		"p",
		8080,
		"Port to start the server on")
}

// initConfig reads in config file and ENV variables if set
func InitConfig() {
	if !strings.HasPrefix(cfgFile, string(filepath.Separator)) {
		cfgFile = fmt.Sprintf("%s%s%s", homeDir, string(filepath.Separator), cfgFile)
	}

	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err != nil {
		helpers.LogError(fmt.Sprintf("Unable to use config file: '%s'. %s", viper.ConfigFileUsed(), err.Error()), "startup - load config")
	}

	repoLocation = helpers.GetRepoPath(repoLocation, homeDir)

	switch helpers.GetRepoTypeString(repoType) {
	case "":
		fallthrough
	case "bolt":
		wlRepo = repository.NewBBoltRepo(repoLocation)
	case "legacy":
		wlRepo = repository.NewYamlFileRepo()
	default:
		helpers.LogWarn(e.RootRepoType, "startup - unknown repo type")
		os.Exit(e.StartupErrors)
	}

	wlConfig = repository.NewYamlConfig(
		filepath.Dir(cfgFile))
	wlService = service.NewWorklogService(wlRepo)
}
