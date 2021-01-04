package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const (
	shortLength = 30
	longLength  = 256
	dateString  = "2000-01-30T0000:00:00Z"
)

func genRandWork() *Work {
	return NewWork(
		helpers.RandString(shortLength),
		helpers.RandString(longLength),
		helpers.RandString(shortLength),
		shortLength,
		[]string{
			helpers.RandString(longLength),
			helpers.RandString(longLength),
		},
		time.Now())
}

func TestNewWork(t *testing.T) {
	validDate, err := time.Parse(time.RFC3339, "1970-12-25T00:00:00Z")
	if err != nil {
		t.Error("Unable to parse initial date")
	}

	var tests = []struct {
		name         string
		wTitle       string
		wDescription string
		wAuthor      string
		wDuration    int
		wTags        []string
		wWhen        time.Time
		expected     *Work
	}{
		{
			name:         "Full work",
			wTitle:       "title",
			wDescription: "description",
			wAuthor:      "who",
			wDuration:    15,
			wTags:        []string{"alpha", "beta"},
			wWhen:        validDate,
			expected: &Work{
				Title:       "title",
				Description: "description",
				Author:      "who",
				Where:       "",
				Duration:    15,
				Tags:        []string{"alpha", "beta"},
				When:        validDate,
			},
		},
		{
			name:         "Unordered tags become ordered",
			wTitle:       "title",
			wDescription: "description",
			wAuthor:      "who",
			wDuration:    15,
			wTags:        []string{"4", "2", "1", "3"},
			wWhen:        validDate,
			expected: &Work{
				Title:       "title",
				Description: "description",
				Author:      "who",
				Where:       "",
				Duration:    15,
				Tags:        []string{"1", "2", "3", "4"},
				When:        validDate,
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := NewWork(
				testItem.wTitle,
				testItem.wDescription,
				testItem.wAuthor,
				testItem.wDuration,
				testItem.wTags,
				testItem.wWhen)
			finished := time.Now()

			// Instead of mocking time.Now(), just set the result of it to the expected value
			testItem.expected.CreatedAt = actual.CreatedAt

			assert.Equal(t, testItem.expected, actual)
			assert.True(t, finished.Add(time.Second*-1).Before(actual.CreatedAt))
		})
	}
}

func TestString(t *testing.T) {
	date, _ := helpers.GetStringAsDateTime(dateString)
	var tests = []struct {
		name string
		work *Work
		exp  string
	}{
		{
			name: "Full work",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], When: %s, CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing title",
			work: &Work{
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Description: %s, Author: %s, Duration: %d, Tags: [%s], When: %s, CreatedAt: %s",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing description",
			work: &Work{
				Title:    "title",
				Author:   "author",
				Duration: 15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Author: %s, Duration: %d, Tags: [%s], When: %s, CreatedAt: %s",
				"title",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing author",
			work: &Work{
				Title:       "title",
				Description: "description",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Duration: %d, Tags: [%s], When: %s, CreatedAt: %s",
				"title",
				"description",
				15,
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing duration",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Tags: [%s], When: %s, CreatedAt: %s",
				"title",
				"description",
				"author",
				"alpha, beta",
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing tags",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				When:        date,
				CreatedAt:   date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, When: %s, CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Empty tags",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags:        []string{},
				When:        date,
				CreatedAt:   date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, When: %s, CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				helpers.TimeFormat(date),
				helpers.TimeFormat(date)),
		}, {
			name: "Missing when",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Basic when",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      time.Time{},
				CreatedAt: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], CreatedAt: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Missing createdAt",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When: date,
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], When: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Basic createdAt",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When:      date,
				CreatedAt: time.Time{},
			},
			exp: fmt.Sprintf("Title: %s, Description: %s, Author: %s, Duration: %d, Tags: [%s], When: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "No fields",
			work: &Work{},
			exp:  "",
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := testItem.work.String()
			assert.Equal(t, testItem.exp, actual)
		})
	}
}

func TestPrettyString(t *testing.T) {
	date, _ := helpers.GetStringAsDateTime(dateString)
	var tests = []struct {
		name string
		work *Work
		exp  string
	}{
		{
			name: "Full work",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When: date,
			},
			exp: fmt.Sprintf("Title: %s\nDescription: %s\nAuthor: %s\nDuration: %d\nTags: [%s]\nWhen: %s",
				"title",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Missing title",
			work: &Work{
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When: date,
			},
			exp: fmt.Sprintf("Description: %s\nAuthor: %s\nDuration: %d\nTags: [%s]\nWhen: %s",
				"description",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Missing description",
			work: &Work{
				Title:    "title",
				Author:   "author",
				Duration: 15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When: date,
			},
			exp: fmt.Sprintf("Title: %s\nAuthor: %s\nDuration: %d\nTags: [%s]\nWhen: %s",
				"title",
				"author",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Missing author",
			work: &Work{
				Title:       "title",
				Description: "description",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When: date,
			},
			exp: fmt.Sprintf("Title: %s\nDescription: %s\nDuration: %d\nTags: [%s]\nWhen: %s",
				"title",
				"description",
				15,
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Missing duration",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Tags: []string{
					"alpha",
					"beta",
				},
				When: date,
			},
			exp: fmt.Sprintf("Title: %s\nDescription: %s\nAuthor: %s\nTags: [%s]\nWhen: %s",
				"title",
				"description",
				"author",
				"alpha, beta",
				helpers.TimeFormat(date)),
		}, {
			name: "Missing tags",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				When:        date,
			},
			exp: fmt.Sprintf("Title: %s\nDescription: %s\nAuthor: %s\nDuration: %d\nWhen: %s",
				"title",
				"description",
				"author",
				15,
				helpers.TimeFormat(date)),
		}, {
			name: "Empty tags",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags:        []string{},
				When:        date,
			},
			exp: fmt.Sprintf("Title: %s\nDescription: %s\nAuthor: %s\nDuration: %d\nWhen: %s",
				"title",
				"description",
				"author",
				15,
				helpers.TimeFormat(date)),
		}, {
			name: "Missing when",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
			},
			exp: fmt.Sprintf("Title: %s\nDescription: %s\nAuthor: %s\nDuration: %d\nTags: [%s]",
				"title",
				"description",
				"author",
				15,
				"alpha, beta"),
		}, {
			name: "Basic when",
			work: &Work{
				Title:       "title",
				Description: "description",
				Author:      "author",
				Duration:    15,
				Tags: []string{
					"alpha",
					"beta",
				},
				When: time.Time{},
			},
			exp: fmt.Sprintf("Title: %s\nDescription: %s\nAuthor: %s\nDuration: %d\nTags: [%s]",
				"title",
				"description",
				"author",
				15,
				"alpha, beta"),
		}, {
			name: "No fields",
			work: &Work{},
			exp:  "",
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := testItem.work.PrettyString()
			assert.Equal(t, testItem.exp, actual)
		})
	}
}

func TestWriteText(t *testing.T) {
	var tests = []struct {
		name   string
		work   *Work
		retErr error
	}{
		{
			name:   "No error",
			work:   genRandWork(),
			retErr: nil,
		}, {
			name:   "Erroring",
			work:   genRandWork(),
			retErr: errors.New(helpers.RandString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			writer := new(MockWriter)
			writer.
				On("Write", []byte(testItem.work.String())).
				Return(1, testItem.retErr)

			actualErr := testItem.work.WriteText(writer)

			writer.AssertExpectations(t)
			writer.AssertCalled(t,
				"Write",
				[]byte(testItem.work.String()),
			)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}

func TestWritePrettyText(t *testing.T) {
	var tests = []struct {
		name   string
		work   *Work
		retErr error
	}{
		{
			name:   "No error",
			work:   genRandWork(),
			retErr: nil,
		}, {
			name:   "Erroring",
			work:   genRandWork(),
			retErr: errors.New(helpers.RandString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			writer := new(MockWriter)
			writer.
				On("Write", []byte(testItem.work.PrettyString())).
				Return(1, testItem.retErr)

			actualErr := testItem.work.WritePrettyText(writer)

			writer.AssertExpectations(t)
			writer.AssertCalled(t,
				"Write",
				[]byte(testItem.work.PrettyString()),
			)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}

func TestWriteAllToPrettyText(t *testing.T) {
	var tests = []struct {
		name   string
		work   []*Work
		retErr error
	}{
		{
			name:   "No error single",
			work:   []*Work{genRandWork()},
			retErr: nil,
		}, {
			name:   "No error double",
			work:   []*Work{genRandWork(), genRandWork()},
			retErr: nil,
		}, {
			name:   "No error quad",
			work:   []*Work{genRandWork(), genRandWork(), genRandWork(), genRandWork()},
			retErr: nil,
		}, {
			name:   "Erroring",
			work:   []*Work{genRandWork()},
			retErr: errors.New(helpers.RandString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			writer := new(MockWriter)

			if testItem.retErr == nil {
				writer.On("Write", []byte("\n")).Return(1, nil)
			}
			for _, element := range testItem.work {
				writer.
					On("Write", []byte(element.PrettyString())).
					Return(1, testItem.retErr)
			}

			actualErr := WriteAllWorkToPrettyText(writer, testItem.work)

			expectedCalled := 1
			if testItem.retErr == nil {
				// for each round, called three more times,
				// except for last round which is only two calls
				expectedCalled = len(testItem.work)*3 - 1
			}

			writer.AssertExpectations(t)
			writer.AssertNumberOfCalls(t, "Write", expectedCalled)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}

func TestWriteYaml(t *testing.T) {
	var tests = []struct {
		name   string
		work   *Work
		retErr error
	}{
		{
			name:   "No error",
			work:   genRandWork(),
			retErr: nil,
		}, {
			name:   "Erroring",
			work:   genRandWork(),
			retErr: errors.New(helpers.RandString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			bytes, _ := yaml.Marshal(testItem.work)
			writer := new(MockWriter)
			writer.
				On("Write", bytes).
				Return(1, testItem.retErr)

			actualErr := testItem.work.WriteYAML(writer)

			writer.AssertExpectations(t)
			writer.AssertCalled(t, "Write", bytes)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}

func TestReadYaml(t *testing.T) {
	// The wonders of .Equals with time's
	date, _ := helpers.GetStringAsDateTime(dateString)
	date, _ = helpers.GetStringAsDateTime(helpers.TimeFormat(date))
	var tests = []struct {
		name    string
		input   string
		expWork *Work
		expErr  error
	}{
		{
			name:  "Full input",
			input: fmt.Sprintf("title: Foo\ndescription: bar\nauthor: possiblellama\nduration: 60\ntags: [1, 2]\nwhen: %s\ncreatedAt: %s", helpers.TimeFormat(date), helpers.TimeFormat(date)),
			expWork: &Work{
				Title:       "Foo",
				Description: "bar",
				Author:      "possiblellama",
				Duration:    60,
				Tags:        []string{"1", "2"},
				When:        date,
				CreatedAt:   date,
			},
			expErr: nil,
		}, {
			name:  "Partial input",
			input: "title: Foo\ndescription: bar\nauthor: possiblellama\nduration: 60",
			expWork: &Work{
				Title:       "Foo",
				Description: "bar",
				Author:      "possiblellama",
				Duration:    60,
				Tags:        []string(nil),
				When:        time.Time{},
				CreatedAt:   time.Time{},
			},
			expErr: nil,
		}, {
			name:    "Invalid fields",
			input:   "foo: bar",
			expWork: &Work{},
			expErr:  nil,
		}, {
			name:    "Invalid format",
			input:   "{\"foo\": \"bar\"}",
			expWork: &Work{},
			expErr:  nil,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actualWork, actualErr := ReadYAML([]byte(testItem.input))

			assert.Equal(t, testItem.expErr, actualErr)
			assert.Equal(t, testItem.expWork, actualWork)
		})
	}
}

func TestWriteJson(t *testing.T) {
	var tests = []struct {
		name   string
		work   *Work
		retErr error
	}{
		{
			name:   "No error",
			work:   genRandWork(),
			retErr: nil,
		}, {
			name:   "Erroring",
			work:   genRandWork(),
			retErr: errors.New(helpers.RandString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			bytes, _ := json.Marshal(testItem.work)
			writer := new(MockWriter)
			writer.
				On("Write", bytes).
				Return(1, testItem.retErr)

			actualErr := testItem.work.WriteJSON(writer)

			writer.AssertExpectations(t)
			writer.AssertCalled(t, "Write", bytes)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}
