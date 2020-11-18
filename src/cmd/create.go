package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new record of work",
	Long: `Creating a new record of work that
the user has created.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().String(
		"title",
		"",
		"A short description of the work done")
	createCmd.Flags().String(
		"description",
		"",
		"A description of the work")
	createCmd.Flags().String(
		"when",
		time.Now().Format(time.RFC3339),
		"When the work was worked")
}
