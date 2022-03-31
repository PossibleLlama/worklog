package helpers

import (
	"testing"
)

func TestLog(t *testing.T) {
	t.Run("Logger", func(t *testing.T) {
		LogDebug("debug message")
		LogInfo("info message")
		LogWarn("warn message")
		LogError("error message")
	})
}
