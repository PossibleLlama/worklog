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
