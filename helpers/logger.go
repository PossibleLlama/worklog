package helpers

import (
	"fmt"
)

func init() {
}

func LogDebug(msg, event string) {
	fmt.Printf("%s\n", msg)
}

func LogInfo(msg, event string) {
	fmt.Printf("%s\n", msg)
}

func LogWarn(msg, event string) {
	fmt.Printf("%s\n", msg)
}

func LogError(msg, event string) {
	fmt.Printf("%s\n", msg)
}
