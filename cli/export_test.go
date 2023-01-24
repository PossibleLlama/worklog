package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/PossibleLlama/worklog/helpers"
	"github.com/PossibleLlama/worklog/service"
	"github.com/stretchr/testify/assert"
)

func setProvidedExportValues(path string) {
	exportPath = path
}

func TestExportArgs(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		assert.Failf(t, err.Error(), fmt.Sprintf("Unable to get home directory during test: %s", err.Error()))
	}
	exportDefaultPath = helpers.RandAlphabeticString(shortLength)

	var tests = []struct {
		name    string
		inPath  string
		expPath string
	}{
		{
			name:    "Variables take defaults",
			inPath:  "",
			expPath: fmt.Sprintf("%s%s%s", homeDir, string(filepath.Separator), exportDefaultPath),
		},
		{
			name:    "Different relative path",
			inPath:  "abc",
			expPath: fmt.Sprintf("%s%sabc", homeDir, string(filepath.Separator)),
		},
		{
			name:    "Different absolute path",
			inPath:  fmt.Sprintf("%sdef", string(filepath.Separator)),
			expPath: fmt.Sprintf("%sdef", string(filepath.Separator)),
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			setProvidedExportValues(testItem.inPath)

			if err := exportArgs(); err != nil {
				assert.Failf(t, err.Error(), "Returned error from exportArgs")
			}

			assert.Equal(t, testItem.expPath, exportPath)
		})
	}
}

func TestExportRun(t *testing.T) {
	var tests = []struct {
		name   string
		path   string
		expErr error
	}{
		{
			name:   "Sends path to service",
			path:   helpers.RandAlphabeticString(shortLength),
			expErr: nil,
		}, {
			name:   "Error propagated",
			path:   helpers.RandAlphabeticString(shortLength),
			expErr: errors.New(helpers.RandAlphabeticString(shortLength)),
		},
	}

	for _, testItem := range tests {
		mockSvc := new(service.MockService)
		mockSvc.On("ExportTo", testItem.path).Return(shortLength, testItem.expErr)

		wlService = mockSvc

		t.Run(testItem.name, func(t *testing.T) {
			setProvidedExportValues(testItem.path)

			actualErr := exportRun()

			mockSvc.AssertExpectations(t)
			mockSvc.AssertCalled(t,
				"ExportTo",
				testItem.path)
			assert.Equal(t, testItem.expErr, actualErr)
		})
	}
}
