package helpers

import (
	"errors"
	"strings"
	"testing"
	"time"
)

const (
	outputTime = "2000-01-02T01:23:00Z"
)

func initilizeTime(t *testing.T, layout, value string) time.Time {
	tm, err := time.Parse(layout, value)
	if err != nil {
		t.Errorf("Initialization of test data failed with %s", err)
	}
	return tm
}

func TestTimeFormat(t *testing.T) {
	rfc3339Time := initilizeTime(t, time.RFC3339, "2000-01-02T01:23:00Z")
	rfc3339NanoTime := initilizeTime(t, time.RFC3339Nano, "2000-01-02T01:23:00.000000000Z")
	rfc1123Time := initilizeTime(t, time.RFC1123, "Mon, 02 Jan 2000 01:23:00 GMT")
	rfc1123ZTime := initilizeTime(t, time.RFC1123Z, "Mon, 02 Jan 2000 01:23:00 +0000")
	rfc822Time := initilizeTime(t, time.RFC822, "02 Jan 00 01:23 GMT")
	rfc822ZTime := initilizeTime(t, time.RFC822Z, "02 Jan 00 01:23 +0000")
	rfc850Time := initilizeTime(t, time.RFC850, "Monday, 02-Jan-00 01:23:00 GMT")

	var tests = []struct {
		name   string
		input  time.Time
		output string
	}{
		{
			name:   "Valid RFC3339 date time is same",
			input:  rfc3339Time,
			output: outputTime,
		}, {
			name:   "Valid RFC3339 nano date time updates",
			input:  rfc3339NanoTime,
			output: outputTime,
		}, {
			name:   "Valid RFC1123 date time updates",
			input:  rfc1123Time,
			output: outputTime,
		}, {
			name:   "Valid RFC1123Z date time updates",
			input:  rfc1123ZTime,
			output: outputTime,
		}, {
			name:   "Valid RFC882 date time updates",
			input:  rfc822Time,
			output: outputTime,
		}, {
			name:   "Valid RFC882Z date time updates",
			input:  rfc822ZTime,
			output: outputTime,
		}, {
			name:   "Valid RFC850 date time updates",
			input:  rfc850Time,
			output: outputTime,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := TimeFormat(testItem.input)

			if actual != testItem.output {
				t.Errorf("Input '%s' was expected to be formated as '%s', but instead '%s'",
					testItem.input,
					testItem.output,
					actual)
			}
		})
	}
}

func TestGetStringAsDateTime(t *testing.T) {
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
			name:             "Valid full date and time with whitespace front",
			dateTimeString:   "\t2000-01-01T09:35:54Z",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid full date and time with whitespace end",
			dateTimeString:   "2000-01-01T09:35:54Z ",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid full date and time with whitespace front and end",
			dateTimeString:   " 2000-01-01T09:35:54Z\n",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Invalid string",
			dateTimeString:   "err",
			dateTimeExpected: time.Now(),
			err:              errors.New("unable to parse string as date"),
		},
		{
			name:             "Valid date with invalid postfix",
			dateTimeString:   "2000-01-01 foo",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              errors.New("cannot parse \" foo\" as \"T\""),
		},
		{
			name:             "Valid date with invalid prefix",
			dateTimeString:   "bar 2000-01-01",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              errors.New("cannot parse \"bar 2000-01-01\" as \"2006\""),
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
