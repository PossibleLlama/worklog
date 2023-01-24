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

const (
	xssHtmlOpen  = "<a href=\"javascript:alert('XSS1')\" onmouseover=\"alert('XSS2')\">"
	xssHtmlClose = "</a>"
)

func genRandWork() *Work {
	return NewWork(
		helpers.RandAlphabeticString(shortLength),
		helpers.RandAlphabeticString(longLength),
		helpers.RandAlphabeticString(shortLength),
		shortLength,
		[]string{
			helpers.RandAlphabeticString(longLength),
			helpers.RandAlphabeticString(longLength),
		},
		time.Now())
}

func TestNewWork(t *testing.T) {
	length := 10
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
				ID:             helpers.RandAlphabeticString(length),
				Revision:       1,
				Title:          "title",
				Description:    "description",
				Author:         "who",
				Duration:       15,
				Tags:           []string{"alpha", "beta"},
				When:           validDate,
				WhenQueryEpoch: validDate.Unix(),
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
				ID:             helpers.RandAlphabeticString(length),
				Revision:       1,
				Title:          "title",
				Description:    "description",
				Author:         "who",
				Duration:       15,
				Tags:           []string{"1", "2", "3", "4"},
				When:           validDate,
				WhenQueryEpoch: validDate.Unix(),
			},
		},
		{
			name:         "Default/empty values are used",
			wTitle:       "",
			wDescription: "",
			wAuthor:      "",
			wDuration:    0,
			wTags:        []string{},
			wWhen:        time.Time{},
			expected: &Work{
				ID:             helpers.RandAlphabeticString(length),
				Revision:       1,
				Title:          "",
				Description:    "",
				Author:         "",
				Duration:       0,
				Tags:           []string{},
				When:           validDate,
				WhenQueryEpoch: validDate.Unix(),
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

			// Instead of mocking ID and time.Now(), just set the
			// result of it to the expected value
			if (testItem.wWhen == time.Time{}) {
				assert.Equal(t, actual.When, actual.CreatedAt)
				assert.Equal(t, actual.WhenQueryEpoch, actual.CreatedAt.Unix())

				testItem.expected.When = actual.When
				testItem.expected.WhenQueryEpoch = actual.When.Unix()
			}
			testItem.expected.ID = actual.ID
			testItem.expected.CreatedAt = actual.CreatedAt

			assert.Equal(t, testItem.expected, actual)
			assert.True(t, finished.Add(time.Second*-1).Before(actual.CreatedAt))
		})
	}
}

func TestUpdate(t *testing.T) {
	wOg := genRandWork()
	wOg.CreatedAt = time.Date(2020, time.January, 30, 23, 59, 0, 0, time.UTC)
	wCopy := Work{
		ID:          wOg.ID,
		Revision:    wOg.Revision,
		Title:       wOg.Title,
		Description: wOg.Description,
		Author:      wOg.Author,
		Duration:    wOg.Duration,
		Tags:        wOg.Tags,
		When:        wOg.When,
		CreatedAt:   wOg.CreatedAt,
	}
	wOg.Update(Work{})

	// Update with no additional fields expected being updated
	assert.Equal(t, wCopy.ID, wOg.ID)
	assert.NotEqual(t, wCopy.Revision, wOg.Revision)
	assert.Equal(t, wCopy.Revision+1, wOg.Revision)
	assert.Equal(t, wCopy.Title, wOg.Title)
	assert.Equal(t, wCopy.Description, wOg.Description)
	assert.Equal(t, wCopy.Author, wOg.Author)
	assert.Equal(t, wCopy.Duration, wOg.Duration)
	assert.Equal(t, wCopy.Tags, wOg.Tags)
	assert.Equal(t, wCopy.When, wOg.When)
	assert.NotEqual(t, wCopy.CreatedAt, wOg.CreatedAt)
	assert.True(t, wCopy.CreatedAt.Before(wOg.CreatedAt))

	new := genRandWork()
	wOg.Update(*new)

	// Update with all new fields
	assert.Equal(t, wCopy.ID, wOg.ID)
	assert.NotEqual(t, wCopy.Revision, wOg.Revision)
	assert.Equal(t, wCopy.Revision+2, wOg.Revision)
	assert.Equal(t, new.Title, wOg.Title)
	assert.Equal(t, new.Description, wOg.Description)
	assert.Equal(t, new.Duration, wOg.Duration)
	assert.Equal(t, new.Author, wOg.Author)
	assert.Equal(t, new.Tags, wOg.Tags)
}

func TestSanitize(t *testing.T) {
	newTag := helpers.RandAlphabeticString(shortLength)
	wSafe := genRandWork()
	wSafe.CreatedAt = time.Date(2020, time.January, 30, 23, 59, 0, 0, time.UTC)
	wXss := Work{
		ID:          wSafe.ID,
		Revision:    wSafe.Revision,
		Title:       xssHtmlOpen + wSafe.Title + xssHtmlClose,
		Description: xssHtmlOpen + wSafe.Description + xssHtmlClose,
		Author:      xssHtmlOpen + wSafe.Author + xssHtmlClose,
		Duration:    wSafe.Duration,
		Tags:        append(wSafe.Tags, xssHtmlOpen+newTag+xssHtmlClose),
		When:        wSafe.When,
		CreatedAt:   wSafe.CreatedAt,
	}
	wXss.Sanitize()

	assert.Equal(t, wSafe.ID, wXss.ID)
	assert.Equal(t, wSafe.Revision, wXss.Revision)
	assert.Equal(t, wSafe.Duration, wXss.Duration)
	assert.Equal(t, wSafe.When, wXss.When)
	assert.Equal(t, wSafe.CreatedAt, wXss.CreatedAt)

	assert.Equal(t, wSafe.Title, wXss.Title)
	assert.Equal(t, wSafe.Description, wXss.Description)
	assert.Equal(t, wSafe.Author, wXss.Author)
	assert.Equal(t, len(wSafe.Tags)+1, len(wXss.Tags))
	assert.Equal(t, append(wSafe.Tags, newTag), wXss.Tags)
}

func TestWorkToPrintWork(t *testing.T) {
	tm := time.Now()
	var tests = []struct {
		name string
		w    Work
		pw   prettyWork
	}{
		{
			name: "Full work",
			w: Work{
				Title:       "Title",
				Description: "Description",
				Author:      "Author",
				Duration:    60,
				Tags:        []string{"1", "2"},
				When:        tm,
				CreatedAt:   time.Now(),
			},
			pw: prettyWork{
				Title:       "Title",
				Description: "Description",
				Author:      "Author",
				Duration:    60,
				Tags:        []string{"1", "2"},
				When:        tm,
			},
		}, {
			name: "Work missing unneeded fields",
			w: Work{
				Title:       "Title",
				Description: "Description",
				Author:      "Author",
				Duration:    60,
				Tags:        []string{"1", "2"},
				When:        tm,
			},
			pw: prettyWork{
				Title:       "Title",
				Description: "Description",
				Author:      "Author",
				Duration:    60,
				Tags:        []string{"1", "2"},
				When:        tm,
			},
		}, {
			name: "Partial work",
			w: Work{
				Title:       "Title",
				Description: "Description",
				Duration:    60,
				When:        tm,
				CreatedAt:   time.Now(),
			},
			pw: prettyWork{
				Title:       "Title",
				Description: "Description",
				Duration:    60,
				When:        tm,
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := workToPrettyWork(testItem.w)

			assert.Equal(t, testItem.pw, actual)
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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			writer := new(mockWriter)
			writer.
				On("Write", []byte(testItem.work.StringNewLine())).
				Return(1, testItem.retErr)

			actualErr := testItem.work.WriteText(writer)

			writer.AssertExpectations(t)
			writer.AssertCalled(t,
				"Write",
				[]byte(testItem.work.StringNewLine()),
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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			writer := new(mockWriter)
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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			writer := new(mockWriter)

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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			bytes, _ := yaml.Marshal(testItem.work)
			writer := new(mockWriter)
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

func TestWritePrettyYaml(t *testing.T) {
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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			bytes, _ := yaml.Marshal(workToPrettyWork(*testItem.work))
			writer := new(mockWriter)
			writer.
				On("Write", bytes).
				Return(1, testItem.retErr)

			actualErr := testItem.work.WritePrettyYAML(writer)

			writer.AssertExpectations(t)
			writer.AssertCalled(t, "Write", bytes)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}

func TestWriteAllToPrettyYaml(t *testing.T) {
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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			wlList := []prettyWork{}
			writer := new(mockWriter)

			for _, element := range testItem.work {
				wlList = append(wlList, workToPrettyWork(*element))
			}
			bytes, _ := yaml.Marshal(wlList)
			writer.On("Write", bytes).Return(1, testItem.retErr)

			actualErr := WriteAllWorkToPrettyYAML(writer, testItem.work)

			writer.AssertExpectations(t)
			writer.AssertNumberOfCalls(t, "Write", 1)
			assert.Equal(t, testItem.retErr, actualErr)
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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			bytes, _ := json.Marshal(testItem.work)
			writer := new(mockWriter)
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

func TestWritePrettyJson(t *testing.T) {
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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			bytes, _ := json.Marshal(workToPrettyWork(*testItem.work))
			writer := new(mockWriter)
			writer.
				On("Write", bytes).
				Return(1, testItem.retErr)

			actualErr := testItem.work.WritePrettyJSON(writer)

			writer.AssertExpectations(t)
			writer.AssertCalled(t, "Write", bytes)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}

func TestWriteAllToPrettyJson(t *testing.T) {
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
			retErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			wlList := []prettyWork{}
			writer := new(mockWriter)

			for _, element := range testItem.work {
				wlList = append(wlList, workToPrettyWork(*element))
			}
			bytes, _ := json.Marshal(wlList)
			writer.On("Write", bytes).Return(1, testItem.retErr)

			actualErr := WriteAllWorkToPrettyJSON(writer, testItem.work)

			writer.AssertExpectations(t)
			writer.AssertNumberOfCalls(t, "Write", 1)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}
