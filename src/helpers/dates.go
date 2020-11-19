package helpers

import (
	"fmt"
	"time"
)

// TimeFormat formats a time to string
func TimeFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}

// GetStringAsDateTime ensures a string is a dateTime
// TODO this should be changed to return an error
func GetStringAsDateTime(element string) (time.Time, error) {
	var dateString string
	if len(element) == 10 {
		dateString = fmt.Sprintf("%sT00:00:00Z", element)
	} else if len(element) == 19 {
		dateString = fmt.Sprintf("%sZ", element)
	} else {
		dateString = element
	}
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return time.Now(), fmt.Errorf("unable to parse string as date. '%s'", err)
	}
	return date, nil
}
