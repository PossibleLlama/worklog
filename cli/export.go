package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/spf13/cobra"
)

var exportPath string
var exportDefaultPath = fmt.Sprintf(".worklog%sexport-%s.json",
	string(filepath.Separator), helpers.TimeFormat(time.Now()))

// exportCmd represents the export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Exports all worklogs",
	Long: `Exports all worklogs in the given repository type
to the given file.`,
	Args: ExportArgs,
	RunE: ExportRun,
}

// ExportArgs public method to validate arguments
func ExportArgs(cmd *cobra.Command, _ []string) error {
	return exportArgs()
}

func exportArgs() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		helpers.LogError(fmt.Sprintf("Unable to get home directory: %s", err.Error()), "export - startup")
		return fmt.Errorf("unable to get home directory: %s", err.Error())
	}

	if exportPath == "" {
		exportPath = exportDefaultPath
	}

	if !strings.HasPrefix(exportPath, string(filepath.Separator)) {
		exportPath = fmt.Sprintf("%s%s%s", homeDir, string(filepath.Separator), exportPath)
	}
	return nil
}

// ExportRun public method to run export
func ExportRun(cmd *cobra.Command, args []string) error {
	return exportRun()
}

func exportRun() error {
	_, err := wlService.ExportTo(exportPath)

	return err
}

func init() {
	rootCmd.AddCommand(exportCmd)

	exportCmd.Flags().StringVar(
		&exportPath,
		"path",
		exportDefaultPath,
		"File path to export results to")
}
