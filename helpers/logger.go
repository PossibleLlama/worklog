package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var NewRelicApp *newrelic.Application

func init() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName("worklog"),
		newrelic.ConfigLicense("NEW_RELIC_LICENSE_KEY"),
		newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		LogDebug(err.Error(), "startup new relic logger")
	}
	NewRelicApp = app
}

func LogDebug(msg, event string) {
	// fmt.Printf("%s\n", msg)
	// logToNewRelic(msg, event, "debug")
}

func LogInfo(msg, event string) {
	fmt.Printf("%s\n", msg)
	logToNewRelic(msg, event, "info")
}

func LogWarn(msg, event string) {
	fmt.Printf("%s\n", msg)
	logToNewRelic(msg, event, "warn")
}

func LogError(msg, event string) {
	fmt.Printf("%s\n", msg)
	logToNewRelic(msg, event, "error")
}

func logToNewRelic(msg, event, level string) {
	if NewRelicApp == nil {
		return
	}
	NewRelicApp.RecordCustomEvent(event, map[string]interface{}{
		"version": Version,
		"time":    time.Now().UTC().Format(time.RFC3339),
		"message": msg,
		"level":   level,
	})
}
