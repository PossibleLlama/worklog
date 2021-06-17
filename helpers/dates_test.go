package helpers

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	outputTimeUTC = "2000-01-02T01:23:00Z"
	outputTime    = "2000-01-02T01:23:00+00:00"
)

func initializeTime(t *testing.T, layout, value string) time.Time {
	tm, err := time.ParseInLocation(layout, value, time.Now().Location())
	if err != nil {
		t.Errorf("Initialization of test data failed with %s", err)
	}
	return tm
}

func initializeTimeBench(b *testing.B, layout, value string) time.Time {
	tm, err := time.ParseInLocation(layout, value, time.Now().Location())
	if err != nil {
		b.Errorf("Initialization of test data failed with %s", err)
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
			input:  initializeTime(t, time.RFC1123, "Sun, 02 Jan 2000 01:23:00 GMT"),
			output: outputTimeUTC,
		}, {
			name:   "Valid RFC1123Z date time updates",
			input:  initializeTime(t, time.RFC1123Z, "Sun, 02 Jan 2000 01:23:00 +0000"),
			output: outputTimeUTC,
		}, {
			name:   "Valid RFC3339 date time is same",
			input:  initializeTime(t, time.RFC3339, "2000-01-02T01:23:00Z"),
			output: outputTimeUTC,
		}, {
			name:   "Valid RFC3339 nano date time updates",
			input:  initializeTime(t, time.RFC3339Nano, "2000-01-02T01:23:00.000000000Z"),
			output: outputTimeUTC,
		}, {
			name:   "Valid RFC850 date time updates",
			input:  initializeTime(t, time.RFC850, "Sunday, 02-Jan-00 01:23:00 GMT"),
			output: outputTimeUTC,
		}, {
			name:   "Valid RFC882 date time updates",
			input:  initializeTime(t, time.RFC822, "02 Jan 00 01:23 GMT"),
			output: outputTimeUTC,
		}, {
			name:   "Valid RFC882Z date time updates",
			input:  initializeTime(t, time.RFC822Z, "02 Jan 00 01:23 +0000"),
			output: outputTimeUTC,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			assert.Equal(t, TimeFormat(testItem.input), testItem.output)
		})
	}
}

func TestGetStringAsDateTime(t *testing.T) {
	expectedDateTimeMidnight := initializeTime(t, time.RFC3339, "2000-01-01T00:00:00+00:00")
	expectedDateTimeMorning := initializeTime(t, time.RFC3339, outputTime)
	expectedDateTimeMorningUTC := initializeTime(t, time.RFC3339, outputTimeUTC)

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
		}, {
			name:             "Valid date. / instead of -",
			dateTimeString:   "2000/01/01",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              nil,
		}, {
			name:             "Valid date. Space instead of -",
			dateTimeString:   "2000 01 01",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              nil,
		}, {
			name:             "Valid date and time",
			dateTimeString:   "2000-01-02T01:23:00",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		}, {
			name:             "Valid date and time. Space instead of T",
			dateTimeString:   "2000-01-02 01:23:00",
			dateTimeExpected: expectedDateTimeMorning,
			err:              nil,
		}, {
			name:             "Valid UTC date and time",
			dateTimeString:   "2000-01-02T01:23:00Z",
			dateTimeExpected: expectedDateTimeMorningUTC,
			err:              nil,
		}, {
			name:             "Valid full date and time with whitespace front",
			dateTimeString:   "\t2000-01-02T01:23:00Z",
			dateTimeExpected: expectedDateTimeMorningUTC,
			err:              nil,
		}, {
			name:             "Valid full date and time with whitespace end",
			dateTimeString:   "2000-01-02T01:23:00Z ",
			dateTimeExpected: expectedDateTimeMorningUTC,
			err:              nil,
		}, {
			name:             "Valid full date and time with whitespace front and end",
			dateTimeString:   " 2000-01-02T01:23:00Z\n",
			dateTimeExpected: expectedDateTimeMorningUTC,
			err:              nil,
		}, {
			name:             "Invalid string",
			dateTimeString:   "err",
			dateTimeExpected: time.Time{},
			err:              errors.New("Could not find format for \"err\""),
		}, {
			name:             "Valid date with invalid postfix",
			dateTimeString:   "2000-01-01 foo",
			dateTimeExpected: expectedDateTimeMidnight,
			err:              nil,
		}, {
			name:             "Valid date with invalid prefix",
			dateTimeString:   "bar 2000-01-01",
			dateTimeExpected: time.Time{},
			err:              errors.New("parsing time \"bar 2000-01-01\" as \"Jan 2000-01-01\": cannot parse \"bar 2000-01-01\" as \"Jan\""),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual, err := GetStringAsDateTime(testItem.dateTimeString)

			if err != nil && testItem.err != nil {
				assert.EqualError(t, err, testItem.err.Error())
			} else if err != nil || testItem.err != nil {
				assert.Fail(t, "Expected or got an error, when the other was not.", fmt.Sprintf("Expected '%s', Actual '%s'", testItem.err, err))
			}

			assert.Equal(t, testItem.dateTimeExpected, actual, fmt.Sprintf("Expected '%s', Actual '%s'", testItem.dateTimeExpected, actual))
		})
	}
}

func BenchmarkGetStringAsDateTime(b *testing.B) {
	var tests = []struct {
		name  string
		input string
	}{
		{
			name:  "YYYY-MM-DD",
			input: "2000-01-02",
		}, {
			name:  "YYYY/MM/DD",
			input: "2000/01/02",
		}, {
			name:  "YYYY MM DD",
			input: "2000 01 02",
		}, {
			name:  "YYYY-MM-DDTHH:mm:SS",
			input: "2000-01-02T01:23:00",
		}, {
			name:  "YYYY-MM-DDTHH:mm:SSZ whitespace start",
			input: "\t2000-01-02T01:23:00Z",
		}, {
			name:  "YYYY-MM-DDTHH:mm:SSZ whitespace end",
			input: "2000-01-02T01:23:00Z\t",
		}, {
			name:  "YYYY-MM-DDTHH:mm:SSZ whitespace both ends",
			input: "\t2000-01-02T01:23:00Z\t",
		},
	}

	for _, testItem := range tests {
		b.Run(testItem.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				if _, err := GetStringAsDateTime(testItem.input); err != nil {
					b.Fail()
				}
			}
		})
	}
}

func TestMidnight(t *testing.T) {
	midnight := initializeTime(t, time.RFC3339, "2000-01-02T00:00:00Z")

	var tests = []struct {
		name       string
		input      time.Time
		is20000102 bool
	}{
		{
			name:       "Second before midnight day",
			input:      initializeTime(t, time.RFC3339, "2000-01-01T23:59:59Z"),
			is20000102: false,
		}, {
			name:       "Midnight stays same",
			input:      initializeTime(t, time.RFC3339, "2000-01-02T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second into day",
			input:      initializeTime(t, time.RFC3339, "2000-01-02T00:00:01Z"),
			is20000102: true,
		}, {
			name:       "Minute into day",
			input:      initializeTime(t, time.RFC3339, "2000-01-02T00:01:00Z"),
			is20000102: true,
		}, {
			name:       "Hour into day",
			input:      initializeTime(t, time.RFC3339, "2000-01-02T01:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second before midday into day",
			input:      initializeTime(t, time.RFC3339, "2000-01-02T11:59:59Z"),
			is20000102: true,
		}, {
			name:       "Midday",
			input:      initializeTime(t, time.RFC3339, "2000-01-02T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second before midnight next day",
			input:      initializeTime(t, time.RFC3339, "2000-01-02T23:59:59Z"),
			is20000102: true,
		}, {
			name:       "Midnight next day",
			input:      initializeTime(t, time.RFC3339, "2000-01-03T00:00:00Z"),
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

func BenchmarkMidnight(b *testing.B) {
	var tests = []struct {
		name  string
		input time.Time
	}{
		{
			name:  "Midnight",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-01T00:00:00Z"),
		}, {
			name:  "Second after",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-01T00:00:01Z"),
		}, {
			name:  "Minute after",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-01T00:01:00Z"),
		}, {
			name:  "Hour after",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-01T01:00:00Z"),
		}, {
			name:  "Second before next day",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-01T23:59:59Z"),
		},
	}

	for _, testItem := range tests {
		b.Run(testItem.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				Midnight(testItem.input)
			}
		})
	}
}

func TestGetPreviousMonday(t *testing.T) {
	monday := initializeTime(t, time.RFC3339, "2000-01-03T00:00:00Z")

	var tests = []struct {
		name       string
		input      time.Time
		is20000102 bool
	}{
		{
			name:       "Second before Monday day",
			input:      initializeTime(t, time.RFC3339, "2000-01-02T23:59:59Z"),
			is20000102: false,
		}, {
			name:       "Midnight Monday stays same",
			input:      initializeTime(t, time.RFC3339, "2000-01-03T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second into Monday",
			input:      initializeTime(t, time.RFC3339, "2000-01-03T00:00:01Z"),
			is20000102: true,
		}, {
			name:       "Minute into Monday",
			input:      initializeTime(t, time.RFC3339, "2000-01-03T00:01:00Z"),
			is20000102: true,
		}, {
			name:       "Hour into Monday",
			input:      initializeTime(t, time.RFC3339, "2000-01-03T01:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second before midday into Monday",
			input:      initializeTime(t, time.RFC3339, "2000-01-03T11:59:59Z"),
			is20000102: true,
		}, {
			name:       "Midday Monday",
			input:      initializeTime(t, time.RFC3339, "2000-01-03T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second before midnight Tuesday",
			input:      initializeTime(t, time.RFC3339, "2000-01-03T23:59:59Z"),
			is20000102: true,
		}, {
			name:       "Midnight Tuesday",
			input:      initializeTime(t, time.RFC3339, "2000-01-04T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midday Tuesday",
			input:      initializeTime(t, time.RFC3339, "2000-01-04T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midnight Wednesday",
			input:      initializeTime(t, time.RFC3339, "2000-01-05T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midday Wednesday",
			input:      initializeTime(t, time.RFC3339, "2000-01-05T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midnight Thursday",
			input:      initializeTime(t, time.RFC3339, "2000-01-06T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midday Thursday",
			input:      initializeTime(t, time.RFC3339, "2000-01-06T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midnight Friday",
			input:      initializeTime(t, time.RFC3339, "2000-01-07T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midday Friday",
			input:      initializeTime(t, time.RFC3339, "2000-01-07T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midnight Saturday",
			input:      initializeTime(t, time.RFC3339, "2000-01-08T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midday Sunday",
			input:      initializeTime(t, time.RFC3339, "2000-01-08T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midnight Sunday",
			input:      initializeTime(t, time.RFC3339, "2000-01-09T00:00:00Z"),
			is20000102: true,
		}, {
			name:       "Midday Sunday",
			input:      initializeTime(t, time.RFC3339, "2000-01-09T12:00:00Z"),
			is20000102: true,
		}, {
			name:       "Second before midnight next Monday",
			input:      initializeTime(t, time.RFC3339, "2000-01-09T23:59:59Z"),
			is20000102: true,
		}, {
			name:       "Midnight next Monday",
			input:      initializeTime(t, time.RFC3339, "2000-01-10T00:00:00Z"),
			is20000102: false,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			if testItem.is20000102 {
				assert.Equal(t, monday, GetPreviousMonday(testItem.input))
			} else {
				assert.NotEqual(t, monday, GetPreviousMonday(testItem.input))
			}
		})
	}
}

func BenchmarkGetPreviousMonday(b *testing.B) {
	var tests = []struct {
		name  string
		input time.Time
	}{
		{
			name:  "Midnight Monday",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-03T00:00:00Z"),
		}, {
			name:  "Midnight Tuesday",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-04T00:00:00Z"),
		}, {
			name:  "Midnight Wednesday",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-05T00:00:00Z"),
		}, {
			name:  "Midnight Thurday",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-06T00:00:00Z"),
		}, {
			name:  "Midnight Friday",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-07T00:00:00Z"),
		}, {
			name:  "Midnight Saturday",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-08T00:00:00Z"),
		}, {
			name:  "Midnight Sunday",
			input: initializeTimeBench(b, time.RFC3339, "2000-01-09T00:00:00Z"),
		},
	}

	for _, testItem := range tests {
		b.Run(testItem.name, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				GetPreviousMonday(testItem.input)
			}
		})
	}
}
