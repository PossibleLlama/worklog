package helpers

import (
	"testing"
)

const (
	event = "testing"
)

func TestLog(t *testing.T) {
	t.Run("Logger", func(t *testing.T) {
		LogDebug("debug message", event)
		LogInfo("info message", event)
		LogWarn("warn message", event)
		LogError("error message", event)
	})
}
