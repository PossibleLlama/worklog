package helpers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var default_repo_path = fmt.Sprintf(".worklog%sworklog.db",
	string(filepath.Separator))

// GetRepoTypeString wrapper for checking viper if not provided
func GetRepoTypeString(arg string) string {
	if arg == "" {
		arg = viper.GetString("repo.type")
	}
	return strings.ToLower(arg)
}

// GetRepoPath wrapper for checking viper if not provided, and falling
// back to a default value if neither are provided.
func GetRepoPath(arg, homeDir string) string {
	var path string
	if arg == "" {
		arg = viper.GetString("repo.path")
		if arg == "" {
			arg = default_repo_path
		}
	}
	// If not absolute path
	if !strings.HasPrefix(arg, string(filepath.Separator)) {
		path = fmt.Sprintf("%s%s%s", homeDir, string(filepath.Separator), arg)
	} else {
		path = arg
	}
	return path
}
