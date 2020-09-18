package main

import (
	"errors"
	"testing"
)

func TestFormatArguments(t *testing.T) {
	emptyFormattedArgs := make(map[string]string)

	var tests = []struct {
		name         string
		rawArgs      []string
		err          error
		returnedArgs map[string]string
	}{
		{
			"No args",
			[]string{},
			errors.New("at least one argument is required"),
			emptyFormattedArgs,
		}, {
			"Single arg",
			[]string{"-r"},
			errors.New("each flag must have an argument"),
			emptyFormattedArgs,
		}, {
			"Single version arg",
			[]string{"--version"},
			nil,
			emptyFormattedArgs,
		}, {
			"Single short version arg",
			[]string{"-v"},
			nil,
			emptyFormattedArgs,
		}, {
			"Single help arg",
			[]string{"--help"},
			nil,
			emptyFormattedArgs,
		}, {
			"Single short help arg",
			[]string{"-h"},
			nil,
			emptyFormattedArgs,
		}, {
			"Both empty args",
			[]string{"", ""},
			errors.New("arguments cannot be empty"),
			emptyFormattedArgs,
		}, {
			"First arg empty",
			[]string{"", "b"},
			errors.New("arguments cannot be empty"),
			emptyFormattedArgs,
		}, {
			"Second arg empty",
			[]string{"a", ""},
			errors.New("arguments cannot be empty"),
			emptyFormattedArgs,
		}, {
			"One pair",
			[]string{"-r", "b"},
			nil,
			map[string]string{"-r": "b"},
		}, {
			"Two pairs",
			[]string{"-r", "b", "-p", "1d"},
			nil,
			map[string]string{"-r": "b", "-p": "1d"},
		}, {
			"Two unordered pairs",
			[]string{"-r", "b", "-p", "1d"},
			nil,
			map[string]string{"-p": "1d", "-r": "b"},
		}, {
			"Duplicate is overridden",
			[]string{"-r", "a", "-r", "b"},
			nil,
			map[string]string{"-r": "b"},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			actualReturnedArgs, actualErr := formatArguments(testItem.rawArgs)

			if testItem.err != nil && actualErr == nil {
				t.Errorf("Should have errored with %s", testItem.err)
			} else if actualErr != nil && testItem.err == nil {
				t.Errorf("Was not supposed to error but did with %s", actualErr)
			} else if testItem.err != nil && testItem.err.Error() != actualErr.Error() {
				t.Errorf("Actual and expected errors are different. Actual: %s Expected: %s", actualErr, testItem.err)
			}
			if len(testItem.returnedArgs) != len(actualReturnedArgs) {
				t.Errorf("Actual and expected number of returned args are different. Actual: %d Expected: %d",
					len(actualReturnedArgs), len(testItem.returnedArgs))
			}
		})
	}
}
