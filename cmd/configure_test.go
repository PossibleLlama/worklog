package cmd

import (
	"errors"
	"strings"
	"testing"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/model"
	"github.com/PossibleLlama/worklog/repository"
	"github.com/stretchr/testify/assert"
)

func setProvidedConfigureValues(author, format string, duration int, rType, rPath string) {
	configProvidedAuthor = author
	configProvidedDuration = duration
	configProvidedFormat = format
	configProvidedRepoType = rType
	configProvidedRepoPath = rPath
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
			setProvidedConfigureValues(
				helpers.RandAlphabeticString(shortLength),
				helpers.RandAlphabeticString(shortLength),
				shortLength,
				helpers.RandAlphabeticString(shortLength),
				helpers.RandAlphabeticString(shortLength))

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
		rType    string
		rPath    string
		expErr   error
	}{
		{
			name:     "Sends model directly to repo",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: shortLength,
			format:   "yaml",
			rType:    "bolt",
			rPath:    "/tmp/foo",
			expErr:   nil,
		}, {
			name:     "Error propogated",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: shortLength,
			format:   "yaml",
			rType:    "bolt",
			rPath:    "/tmp/foo",
			expErr:   errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		cfg := model.NewConfig(model.Defaults{
			Author:   testItem.author,
			Format:   testItem.format,
			Duration: testItem.duration,
		}, model.Repo{
			Type: testItem.rType,
			Path: testItem.rPath,
		})

		mockRepo := new(repository.MockRepo)
		mockRepo.On("SaveConfig", cfg).Return(testItem.expErr)
		wlConfig = mockRepo

		t.Run(testItem.name, func(t *testing.T) {
			setProvidedConfigureValues(
				testItem.author,
				testItem.format,
				testItem.duration,
				testItem.rType,
				testItem.rPath)

			actualErr := configRun()

			mockRepo.AssertExpectations(t)
			mockRepo.AssertCalled(t,
				"SaveConfig",
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
		rType    string
		rPath    string
		expErr   error
	}{
		{
			name:     "Override all variables",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: shortLength,
			format:   "yaml",
			rType:    "bolt",
			rPath:    "/tmp/foo",
			expErr:   nil,
		}, {
			name:     "Empty arguments throws error",
			author:   "",
			duration: -1,
			format:   "",
			rType:    "",
			rPath:    "",
			expErr:   errors.New("overrideDefaults requires at least one argument"),
		}, {
			name:     "Zero duration",
			author:   "",
			duration: 0,
			format:   "",
			rType:    "",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "Non empty author",
			author:   helpers.RandAlphabeticString(shortLength),
			duration: -1,
			format:   "",
			rType:    "",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "Padded author",
			author:   " " + helpers.RandAlphabeticString(shortLength),
			duration: -1,
			format:   "",
			rType:    "",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "Blank author",
			author:   " ",
			duration: -1,
			format:   "",
			rType:    "",
			rPath:    "",
			expErr:   errors.New("overrideDefaults requires at least one argument"),
		}, {
			name:     "yaml format",
			author:   "",
			duration: -1,
			format:   "yaml",
			rType:    "",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "padded yaml format",
			author:   "",
			duration: -1,
			format:   "\tyaml\n",
			rType:    "",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "Pretty format",
			author:   "",
			duration: -1,
			format:   "pretty",
			rType:    "",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "json format",
			author:   "",
			duration: -1,
			format:   "json",
			rType:    "",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "bolt repo type",
			author:   "",
			duration: -1,
			format:   "",
			rType:    "bolt",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "legacy repo type",
			author:   "",
			duration: -1,
			format:   "",
			rType:    "legacy",
			rPath:    "",
			expErr:   nil,
		}, {
			name:     "absolute repo path",
			author:   "",
			duration: -1,
			format:   "",
			rType:    "",
			rPath:    "/tmp/foo",
			expErr:   nil,
		}, {
			name:     "relative repo path",
			author:   "",
			duration: -1,
			format:   "",
			rType:    "",
			rPath:    "./foo",
			expErr:   nil,
		}, {
			name:     "Blank format",
			author:   "",
			duration: -1,
			format:   " ",
			rType:    "",
			rPath:    "",
			expErr:   errors.New("overrideDefaults requires at least one argument"),
		}, {
			name:     "Invalid format",
			author:   "",
			duration: -1,
			format:   helpers.RandAlphabeticString(shortLength),
			rType:    "",
			rPath:    "",
			expErr:   errors.New("format is not valid"),
		}, {
			name:     "Blank repo type",
			author:   "",
			duration: -1,
			format:   "",
			rType:    " ",
			rPath:    "",
			expErr:   errors.New("overrideDefaults requires at least one argument"),
		}, {
			name:     "Blank repo path",
			author:   "",
			duration: -1,
			format:   "",
			rType:    "",
			rPath:    " ",
			expErr:   errors.New("overrideDefaults requires at least one argument"),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedConfigureValues(
				testItem.author,
				testItem.format,
				testItem.duration,
				testItem.rType,
				testItem.rPath)

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
