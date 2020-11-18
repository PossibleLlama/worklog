package cmd

import (
	"fmt"

	"github.com/PossibleLlama/worklog/helpers"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of the app",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(helpers.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
