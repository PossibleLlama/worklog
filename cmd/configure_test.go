package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/service"
	"github.com/stretchr/testify/assert"
)

func setProvidedConfigureValues(author, format string, duration int) {
	configProvidedAuthor = author
	configProvidedDuration = duration
	configProvidedFormat = format
}

func TestConfigArgs(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{
			name: "Variables take defaults",
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedConfigureValues(helpers.RandAlphabeticString(shortLength), helpers.RandAlphabeticString(shortLength), shortLength)

			if err := configArgs(); err != nil {
				assert.Failf(t, err.Error(), "Returned error from configArgs")
			}

			assert.Equal(t, configDefaultAuthor, configProvidedAuthor)
			assert.Equal(t, configDefaultDuration, configProvidedDuration)
			assert.Equal(t, configDefaultFormat, configProvidedFormat)
		})
	}
}

func TestConfigRun(t *testing.T) {
	var tests = []struct {
		name     string
		author   string
		duration int
		format   string
		expErr   error
	}{
		{
			name:     "Sends model to service",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: shortLength,
			format:   "yaml",
			expErr:   nil,
		}, {
			name:     "Error propogated",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: shortLength,
			format:   "yaml",
			expErr:   errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		cfg := model.NewConfig(testItem.author, testItem.format, testItem.duration)

		mockService := new(service.MockService)
		mockService.On("Configure", cfg).Return(testItem.expErr)
		wlService = mockService

		t.Run(testItem.name, func(t *testing.T) {
			setProvidedConfigureValues(testItem.author, testItem.format, testItem.duration)

			actualErr := configRun()

			mockService.AssertExpectations(t)
			mockService.AssertCalled(t,
				"Configure",
				cfg)
			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}

func TestOverrideDefaultsArgs(t *testing.T) {
	var tests = []struct {
		name     string
		author   string
		duration int
		format   string
		expErr   error
	}{
		{
			name:     "Override all variables",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: shortLength,
			format:   "yaml",
			expErr:   nil,
		}, {
			name:     "Empty arguments throws error",
			author:   "",
			duration: -1,
			format:   "",
			expErr:   errors.New("overrideDefaults requires at least one argument"),
		}, {
			name:     "Zero duration",
			author:   "",
			duration: 0,
			format:   "",
			expErr:   nil,
		}, {
			name:     "Non empty author",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: -1,
			format:   "",
			expErr:   nil,
		}, {
			name:     "Padded author",
			author:   " " + helpers.RandAlphabeticString(shortLength),
			duration: -1,
			format:   "",
			expErr:   nil,
		}, {
			name:     "Blank author",
			author:   " ",
			duration: -1,
			format:   "",
			expErr:   errors.New("overrideDefaults requires at least one argument"),
		}, {
			name:     "yaml format",
			author:   "",
			duration: -1,
			format:   "yaml",
			expErr:   nil,
		}, {
			name:     "padded yaml format",
			author:   "",
			duration: -1,
			format:   "\tyaml\n",
			expErr:   nil,
		}, {
			name:     "Pretty format",
			author:   "",
			duration: -1,
			format:   "pretty",
			expErr:   nil,
		}, {
			name:     "json format",
			author:   "",
			duration: -1,
			format:   "json",
			expErr:   nil,
		}, {
			name:     "Blank format",
			author:   "",
			duration: -1,
			format:   " ",
			expErr:   errors.New("overrideDefaults requires at least one argument"),
		}, {
			name:     "Invalid format",
			author:   "",
			duration: -1,
			format:   helpers.RandAlphabeticString(shortLength),
			expErr:   errors.New("format is not valid"),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedConfigureValues(testItem.author, testItem.format, testItem.duration)

			actualErr := overrideDefaultsArgs()

			assert.Equal(t, strings.TrimSpace(testItem.author), configProvidedAuthor)
			assert.Equal(t, strings.TrimSpace(testItem.format), configProvidedFormat)
			if testItem.expErr == nil {
				assert.Nil(t, actualErr)
				if testItem.duration >= 0 {
					assert.Equal(t, testItem.duration, configProvidedDuration)
				} else {
					assert.Equal(t, configDefaultDuration, configProvidedDuration)
				}
			} else {
				assert.EqualError(t, testItem.expErr, actualErr.Error())
			}
		})
	}
}
