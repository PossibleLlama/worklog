package helpers

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandStringLength(t *testing.T) {
	var tests = []struct {
		name   string
		strLen int
		expLen int
	}{
		{
			name:   "0 length",
			strLen: 0,
			expLen: 0,
		}, {
			name:   "1 length",
			strLen: 1,
			expLen: 1,
		}, {
			name:   "Negative length",
			strLen: -1,
			expLen: 0,
		}, {
			name:   "100 length",
			strLen: 100,
			expLen: 100,
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			assert.Len(t, RandAlphabeticString(testItem.strLen), testItem.expLen)
		})
	}
}

func TestRandAlphabeticStringRandomness(t *testing.T) {
	max := int(math.Pow(2, 10))
	seen := make([]string, 0)

	for i := 0; i < max; i++ {
		str := RandAlphabeticString(10)
		assert.NotSubset(t, seen, []string{str})

		seen = append(seen, str)
	}
}

func TestRandHexadecimalStringRandomness(t *testing.T) {
	max := int(math.Pow(2, 10))
	seen := make([]string, 0)

	for i := 0; i < max; i++ {
		str := RandHexAlphaNumericString(10)
		assert.NotSubset(t, seen, []string{str})

		seen = append(seen, str)
	}
}

func TestDeduplicateString(t *testing.T) {
	var tests = []struct {
		name string
		in   []string
		exp  []string
	}{
		{
			name: "Empty",
			in:   []string{},
			exp:  []string{},
		}, {
			name: "1 item",
			in:   []string{"a"},
			exp:  []string{"a"},
		}, {
			name: "2 items",
			in:   []string{"a", "b"},
			exp:  []string{"a", "b"},
		}, {
			name: "1 item repeated",
			in:   []string{"a", "a"},
			exp:  []string{"a"},
		}, {
			name: "2 items repeated",
			in:   []string{"a", "a", "b", "b"},
			exp:  []string{"a", "b"},
		}, {
			name: "3 items, some repeated",
			in:   []string{"a", "a", "b", "c", "b"},
			exp:  []string{"a", "b", "c"},
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.name, func(t *testing.T) {
			assert.Equal(t, len(DeduplicateString(testItem.in)), len(testItem.exp))
			for index, element := range DeduplicateString(testItem.in) {
				assert.Equal(t, testItem.exp[index], element)
			}
		})
	}
}
