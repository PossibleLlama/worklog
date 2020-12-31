package helpers

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	var tests = []struct {
		name   string
		input  time.Time
		output string
	}{
		{
			name:   "Valid RFC1123 date time updates",
			input:  initilizeTime(t, time.RFC1123, "Mon, 02 Jan 2000 01:23:00 GMT"),
			output: outputTime,
		}, {
			name:   "Valid RFC1123Z date time updates",
			input:  initilizeTime(t, time.RFC1123Z, "Mon, 02 Jan 2000 01:23:00 +0000"),
			output: outputTime,
		}, {
			name:   "Valid RFC3339 date time is same",
			input:  initilizeTime(t, time.RFC3339, "2000-01-02T01:23:00Z"),
			output: outputTime,
		}, {
			name:   "Valid RFC3339 nano date time updates",
			input:  initilizeTime(t, time.RFC3339Nano, "2000-01-02T01:23:00.000000000Z"),
			output: outputTime,
		}, {
			name:   "Valid RFC850 date time updates",
			input:  initilizeTime(t, time.RFC850, "Monday, 02-Jan-00 01:23:00 GMT"),
			output: outputTime,
		}, {
			name:   "Valid RFC882 date time updates",
			input:  initilizeTime(t, time.RFC822, "02 Jan 00 01:23 GMT"),
			output: outputTime,
		}, {
			name:   "Valid RFC882Z date time updates",
			input:  initilizeTime(t, time.RFC822Z, "02 Jan 00 01:23 +0000"),
			output: outputTime,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			assert.Equal(t, TimeFormat(testItem.input), testItem.output)
		})
	}
}

func TestGetStringAsDateTime(t *testing.T) {
	expectedDateTimeMidnight := initilizeTime(t, time.RFC3339, "2000-01-01T00:00:00Z")
	expectedDateTimeMorning := initilizeTime(t, time.RFC3339, outputTime)

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
			dateTimeString:   "2000-01-02T01:23:00",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid date and time. Space instead of T",
			dateTimeString:   "2000-01-02 01:23:00",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid full date and time",
			dateTimeString:   "2000-01-02T01:23:00Z",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid full date and time with whitespace front",
			dateTimeString:   "\t2000-01-02T01:23:00Z",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid full date and time with whitespace end",
			dateTimeString:   "2000-01-02T01:23:00Z ",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Valid full date and time with whitespace front and end",
			dateTimeString:   " 2000-01-02T01:23:00Z\n",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		},
		{
			name:             "Invalid string",
			dateTimeString:   "err",
			dateTimeExpected: time.Now(),
			err:              errors.New("unable to parse string as date. 'parsing time \"err\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"err\" as \"2006\"'"),
		},
		{
			name:             "Valid date with invalid postfix",
			dateTimeString:   "2000-01-01 foo",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              errors.New("unable to parse string as date. 'parsing time \"2000-01-01 foo\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \" foo\" as \"T\"'"),
		},
		{
			name:             "Valid date with invalid prefix",
			dateTimeString:   "bar 2000-01-01",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              errors.New("unable to parse string as date. 'parsing time \"bar 2000-01-01\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"bar 2000-01-01\" as \"2006\"'"),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual, err := GetStringAsDateTime(testItem.dateTimeString)

			if err != nil {
				assert.EqualError(t, testItem.err, err.Error())
			} else if testItem.err != nil {
				assert.EqualError(t, err, testItem.err.Error())
			} else {
				assert.Equal(t, testItem.dateTimeExpected, actual)
			}
		})
	}
}

func TestMidnight(t *testing.T) {
	midnight := initilizeTime(t, time.RFC3339, "2000-01-02T00:00:00Z")

	var tests = []struct {
		name       string
		input      time.Time
		is20000102 bool
	}{
		{
			name:       "Second before midnight day",
			input:      initilizeTime(t, time.RFC3339, "2000-01-01T23:59:59Z"),
			is20000102: false,
		}, {
			name:       "Midnight stays same",
			input:      initilizeTime(t, time.RFC3339, "2000-01-02T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second into day",
			input:      initilizeTime(t, time.RFC3339, "2000-01-02T00:00:01Z"),
			is20000102: true,
		}, {
			name:       "Minute into day",
			input:      initilizeTime(t, time.RFC3339, "2000-01-02T00:01:00Z"),
			is20000102: true,
		}, {
			name:       "Hour into day",
			input:      initilizeTime(t, time.RFC3339, "2000-01-02T01:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second before midday into day",
			input:      initilizeTime(t, time.RFC3339, "2000-01-02T11:59:59Z"),
			is20000102: true,
		}, {
			name:       "Midday",
			input:      initilizeTime(t, time.RFC3339, "2000-01-02T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second before midnight next day",
			input:      initilizeTime(t, time.RFC3339, "2000-01-02T23:59:59Z"),
			is20000102: true,
		}, {
			name:       "Midnight next day",
			input:      initilizeTime(t, time.RFC3339, "2000-01-03T00:00:00Z"),
			is20000102: false,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			if testItem.is20000102 {
				assert.Equal(t, midnight, Midnight(testItem.input))
			} else {
				assert.NotEqual(t, midnight, Midnight(testItem.input))
			}
		})
	}
}
