package helpers

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	expectedDateTimeMidnight, err := time.Parse(time.RFC3339, "2000-01-01T00:00:00Z")
	if err != nil {
		t.Errorf("Initialization of test data failed with %s", err)
	}
	expectedDateTimeMorning, err := time.Parse(time.RFC3339, "2000-01-01T09:35:54Z")
	if err != nil {
		t.Errorf("Initialization of test data failed with %s", err)
	}

	var tests = []struct {
		name             string
		dateTimeString   string
		dateTimeExpected time.Time
		err              error
	}{
		{
			name:             "Valid date",
			dateTimeString:   "2000-01-01",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              nil,
		},
		{
			name:             "Valid date. / instead of -",
			dateTimeString:   "2000/01/01",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              nil,
		},
		{
			name:             "Valid date. Space instead of -",
			dateTimeString:   "2000 01 01",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              nil,
		},
		{
			name:             "Valid date and time",
			dateTimeString:   "2000-01-01T09:35:54",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid date and time. Space instead of T",
			dateTimeString:   "2000-01-01 09:35:54",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid full date and time",
			dateTimeString:   "2000-01-01T09:35:54Z",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Invalid string",
			dateTimeString:   "err",
			dateTimeExpected: time.Now(),
			err:              errors.New("unable to parse string as date"),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual, err := GetStringAsDateTime(testItem.dateTimeString)

			if err != nil || testItem.err != nil {
				if testItem.err == nil {
					t.Errorf("Produced error %s when none expected", err)
				} else if err == nil {
					t.Error("Expected error to be produced, but none returned")
				} else if !strings.Contains(err.Error(), testItem.err.Error()) {
					t.Errorf("Expected error to contain '%s', but was '%s'",
						testItem.err.Error(), err.Error())
				}
			} else {
				if !actual.Equal(testItem.dateTimeExpected) {
					t.Errorf("Expected %s to be returned from parsing %s, instead got %s", testItem.dateTimeExpected, testItem.dateTimeString, actual)
				}
			}
		})
	}
}