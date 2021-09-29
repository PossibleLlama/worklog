package model

import (
	"testing"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

const path = "/tmp/foo"

func TestNewConfig(t *testing.T) {
	var tests = []struct {
		name     string
		author   string
		format   string
		duration int
		rType    string
		rPath    string
		expected *Config
	}{
		{
			name:     "Full config: pretty",
			author:   "Author",
			format:   "pretty",
			duration: 60,
			rType:    "bolt",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "pretty",
					Duration: 60,
				},
				Repo: Repo{
					Type: "bolt",
					Path: path,
				},
			},
		}, {
			name:     "Full config: yaml",
			author:   "Author",
			format:   "yaml",
			duration: 60,
			rType:    "bolt",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yaml",
					Duration: 60,
				},
				Repo: Repo{
					Type: "bolt",
					Path: path,
				},
			},
		}, {
			name:     "Full config: yml",
			author:   "Author",
			format:   "yml",
			duration: 60,
			rType:    "bolt",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yml",
					Duration: 60,
				},
				Repo: Repo{
					Type: "bolt",
					Path: path,
				},
			},
		}, {
			name:     "Full config: json",
			author:   "Author",
			format:   "json",
			duration: 60,
			rType:    "bolt",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "json",
					Duration: 60,
				},
				Repo: Repo{
					Type: "bolt",
					Path: path,
				},
			},
		}, {
			name:     "Full config: invalid format",
			author:   "Author",
			format:   "foo",
			duration: 60,
			rType:    "bolt",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "",
					Duration: 60,
				},
				Repo: Repo{
					Type: "bolt",
					Path: path,
				},
			},
		}, {
			name:     "Missing format",
			author:   "Author",
			format:   "",
			duration: 60,
			rType:    "bolt",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "",
					Duration: 60,
				},
				Repo: Repo{
					Type: "bolt",
					Path: path,
				},
			},
		}, {
			name:     "Missing author",
			author:   "",
			format:   "yaml",
			duration: 60,
			rType:    "bolt",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "",
					Format:   "yaml",
					Duration: 60,
				},
				Repo: Repo{
					Type: "bolt",
					Path: path,
				},
			},
		}, {
			name:     "Missing repo type",
			author:   "Author",
			format:   "yaml",
			duration: 60,
			rType:    "",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yaml",
					Duration: 60,
				},
				Repo: Repo{
					Type: "",
					Path: path,
				},
			},
		}, {
			name:     "Missing repo path",
			author:   "Author",
			format:   "yaml",
			duration: 60,
			rType:    "bolt",
			rPath:    "",
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yaml",
					Duration: 60,
				},
				Repo: Repo{
					Type: "bolt",
					Path: "",
				},
			},
		}, {
			name:     "Zero duration",
			author:   "Author",
			format:   "yaml",
			duration: 0,
			rType:    "bolt",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yaml",
					Duration: 0,
				},
				Repo: Repo{
					Type: "bolt",
					Path: path,
				},
			},
		}, {
			name:     "Invalid repo type",
			author:   "Author",
			format:   "yaml",
			duration: 60,
			rType:    "invalid",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yaml",
					Duration: 60,
				},
				Repo: Repo{
					Type: "",
					Path: path,
				},
			},
		}, {
			name:     "Just author",
			author:   "Author",
			format:   "",
			duration: 0,
			rType:    "",
			rPath:    "",
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "",
					Duration: 0,
				},
				Repo: Repo{
					Type: "",
					Path: "",
				},
			},
		}, {
			name:     "Just format",
			author:   "",
			format:   "yaml",
			duration: 0,
			rType:    "",
			rPath:    "",
			expected: &Config{
				Defaults: Defaults{
					Author:   "",
					Format:   "yaml",
					Duration: 0,
				},
				Repo: Repo{
					Type: "",
					Path: "",
				},
			},
		}, {
			name:     "Just duration",
			author:   "",
			format:   "",
			duration: 60,
			rType:    "",
			rPath:    "",
			expected: &Config{
				Defaults: Defaults{
					Author:   "",
					Format:   "",
					Duration: 60,
				},
				Repo: Repo{
					Type: "",
					Path: "",
				},
			},
		}, {
			name:     "Just repo type",
			author:   "",
			format:   "",
			duration: 0,
			rType:    "bolt",
			rPath:    "",
			expected: &Config{
				Defaults: Defaults{
					Author:   "",
					Format:   "",
					Duration: 0,
				},
				Repo: Repo{
					Type: "bolt",
					Path: "",
				},
			},
		}, {
			name:     "Just repo path",
			author:   "",
			format:   "",
			duration: 0,
			rType:    "",
			rPath:    path,
			expected: &Config{
				Defaults: Defaults{
					Author:   "",
					Format:   "",
					Duration: 0,
				},
				Repo: Repo{
					Type: "",
					Path: path,
				},
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := NewConfig(Defaults{
				Author:   testItem.author,
				Format:   testItem.format,
				Duration: testItem.duration,
			},
				Repo{
					Type: testItem.rType,
					Path: testItem.rPath,
				})

			assert.Equal(t, testItem.expected, actual)
		})
	}
}

func TestConfigWriteYaml(t *testing.T) {
	var tests = []struct {
		name   string
		cfg    *Config
		retErr error
	}{
		{
			name: "Full config",
			cfg: &Config{
				Defaults: Defaults{
					Author:   helpers.RandAlphabeticString(shortLength),
					Format:   helpers.RandAlphabeticString(shortLength),
					Duration: shortLength,
				},
			},
			retErr: nil,
		}, {
			name: "Partial config",
			cfg: &Config{
				Defaults: Defaults{
					Author: helpers.RandAlphabeticString(shortLength),
					Format: helpers.RandAlphabeticString(shortLength),
				},
			},
			retErr: nil,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			bytes, _ := yaml.Marshal(testItem.cfg)
			writer := new(mockWriter)
			writer.
				On("Write", bytes).
				Return(1, testItem.retErr)

			actualErr := testItem.cfg.WriteYAML(writer)

			writer.AssertExpectations(t)
			writer.AssertCalled(t, "Write", bytes)
			assert.Equal(t, testItem.retErr, actualErr)
		})
	}
}
