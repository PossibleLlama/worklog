package helpers

import (
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

// These regex'es are not expected to be exhaustive.
// They will be used to check whether this is likely
// to be a date, and not a random string
const dateRegex = `[0-9]{4}[-/ ][0-1][0-9][-/ ][0-3][0-9]`

// TimeFormat formats a time to string
func TimeFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}

// GetStringAsDateTime ensures a string is a dateTime
func GetStringAsDateTime(rawElement string) (time.Time, error) {
	element := strings.TrimSpace(rawElement)

	isDate, dateErr := regexp.MatchString("^"+dateRegex, element)
	if dateErr != nil {
		return time.Time{}, dateErr
	}
	if isDate {
		element = replaceAtIndex(element, '-', 4)
		element = replaceAtIndex(element, '-', 7)
	}

	tm, err := dateparse.ParseLocal(element)
	if err != nil {
		return time.Time{}, err
	}
	return tm, nil
}

// Midnight tonight
func Midnight(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// GetPreviousMonday getting the most recent Monday
func GetPreviousMonday(originalTime time.Time) time.Time {
	t := Midnight(originalTime)
	for i := 0; i <= 6; i++ {
		if t.Weekday() == time.Monday {
			return t
		}
		t = t.AddDate(0, 0, -1)
	}
	return originalTime
}
