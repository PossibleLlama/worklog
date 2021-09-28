package helpers

import (
	"strings"

	"github.com/spf13/viper"
)

// GetRepoTypeString wrapper for checking viper if not provided
func GetRepoTypeString(arg string) string {
	if arg == "" {
		arg = viper.GetString("repo")
	}
	return strings.ToLower(arg)
}
