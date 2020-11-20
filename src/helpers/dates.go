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

// TimeFormat formats a time to string
func TimeFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}

// GetStringAsDateTime ensures a string is a dateTime
func GetStringAsDateTime(element string) (time.Time, error) {
	var dateString string

	isDate, dateErr := regexp.MatchString(dateRegex, element)
	isDateTime, dateTimeErr := regexp.MatchString(dateTimeRegex, element)

	if dateErr != nil || dateTimeErr != nil {
		return time.Now(), fmt.Errorf("unable to parse string as date")
	}

	if isDateTime {
		dateString = fmt.Sprintf("%s-%s-%sT%s:%s:%sZ",
			string(element[0:4]),
			string(element[5:7]),
			string(element[8:10]),
			string(element[11:13]),
			string(element[14:16]),
			string(element[17:19]))
	} else if isDate {
		dateString = fmt.Sprintf("%s-%s-%sT00:00:00Z",
			string(element[0:4]),
			string(element[5:7]),
			string(element[8:10]))
	} else {
		dateString = element
	}
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return time.Now(), fmt.Errorf("unable to parse string as date. '%s'", err)
	}
	return date, nil
}