package model

import (
	"testing"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestNewConfig(t *testing.T) {
	var tests = []struct {
		name     string
		author   string
		format   string
		duration int
		expected *Config
	}{
		{
			name:     "Full config: pretty",
			author:   "Author",
			format:   "pretty",
			duration: 60,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "pretty",
					Duration: 60,
				},
			},
		}, {
			name:     "Full config: yaml",
			author:   "Author",
			format:   "yaml",
			duration: 60,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yaml",
					Duration: 60,
				},
			},
		}, {
			name:     "Full config: yml",
			author:   "Author",
			format:   "yml",
			duration: 60,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yml",
					Duration: 60,
				},
			},
		}, {
			name:     "Full config: json",
			author:   "Author",
			format:   "json",
			duration: 60,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "json",
					Duration: 60,
				},
			},
		}, {
			name:     "Full config: invalid format",
			author:   "Author",
			format:   "foo",
			duration: 60,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "",
					Duration: 60,
				},
			},
		}, {
			name:     "Missing format",
			author:   "Author",
			format:   "",
			duration: 60,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "",
					Duration: 60,
				},
			},
		}, {
			name:     "Missing author",
			author:   "",
			format:   "yaml",
			duration: 60,
			expected: &Config{
				Defaults: Defaults{
					Author:   "",
					Format:   "yaml",
					Duration: 60,
				},
			},
		}, {
			name:     "Zero duration",
			author:   "Author",
			format:   "yaml",
			duration: 0,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "yaml",
					Duration: 0,
				},
			},
		}, {
			name:     "Just author",
			author:   "Author",
			format:   "",
			duration: 0,
			expected: &Config{
				Defaults: Defaults{
					Author:   "Author",
					Format:   "",
					Duration: 0,
				},
			},
		}, {
			name:     "Just format",
			author:   "",
			format:   "yaml",
			duration: 0,
			expected: &Config{
				Defaults: Defaults{
					Author:   "",
					Format:   "yaml",
					Duration: 0,
				},
			},
		}, {
			name:     "Just duration",
			author:   "",
			format:   "",
			duration: 60,
			expected: &Config{
				Defaults: Defaults{
					Author:   "",
					Format:   "",
					Duration: 60,
				},
			},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actual := NewConfig(testItem.author, testItem.format, testItem.duration)

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
