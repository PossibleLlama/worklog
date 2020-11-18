package helpers

import (
	"fmt"
	"os"
	"time"
)

// TimeFormat formats a time to string
func TimeFormat(t time.Time) string {
	return t.Format(time.RFC3339)
}

// GetStringAsDateTime ensures a string is a dateTime
// TODO this should be changed to return an error
func GetStringAsDateTime(element string) time.Time {
	var dateString string
	if len(element) == 10 {
		dateString = fmt.Sprintf("%sT00:00:00Z", element)
	} else {
		dateString = element
	}
	date, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		fmt.Printf("Date to print from is not a valid date. %s\n", err.Error())
		os.Exit(1)
	}
	return date
}
