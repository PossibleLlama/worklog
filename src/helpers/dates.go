package helpers

import (
	"fmt"
	"regexp"
	"time"
)

// These regex'es are not expected to be exhaustive.
// They will be used to check whether this is likely
// to be a date, and not a random string
const dateRegex = `[0-9]{4}[-/ ][0-1][0-9][-/ ][0-3][0-9]`
const timeRegex = `[0-2][0-9]:[0-5][0-9]:[0-5][0-9]`
const dateTimeRegex = dateRegex + `[\sT]` + timeRegex
const dateTimeRFC3339Regex = dateTimeRegex + `Z`

// TimeFormat formats a time to string
func TimeFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}

// GetStringAsDateTime ensures a string is a dateTime
func GetStringAsDateTime(element string) (time.Time, error) {
	var dateString string

	_, dateErr := regexp.MatchString(dateRegex, element)
	isDateTime, dateTimeErr := regexp.MatchString(dateTimeRegex, element)
	isRFC339, rfc399Err := regexp.MatchString(dateTimeRFC3339Regex, element)

	if dateErr != nil || dateTimeErr != nil || rfc399Err != nil {
		return time.Now(), fmt.Errorf("unable to parse string as date")
	}

	if isRFC339 {
		dateString = element
	} else if isDateTime {
		dateString = fmt.Sprintf("%sZ", element)
	} else {
		dateString = fmt.Sprintf("%sT00:00:00Z", element)
	}
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return time.Now(), fmt.Errorf("unable to parse string as date. '%s'", err)
	}
	return date, nil
}
