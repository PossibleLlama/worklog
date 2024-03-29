package helpers

import (
	"crypto/rand"
	"log"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

const (
	alphabeticBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hexadecimalBytes = "abcdef1234567890"
	letterIdxBits    = 6
	letterIdxMask    = 1<<letterIdxBits - 1
)

// RegexCaseInsensitive regex to be prepended to not care about case
// of characters
const RegexCaseInsensitive = "(?i)"

var policy = bluemonday.StrictPolicy()

// RandAlphabeticString Generator function of a random series of characters
// Uses a-zA-Z character set
func RandAlphabeticString(n int) string {
	return randString(n, alphabeticBytes)
}

// RandHexAlphaNumericString Generator function of a random series of characters
// Uses a-f0-9 character set
func RandHexAlphaNumericString(n int) string {
	return randString(n, hexadecimalBytes)
}

// From https://stackoverflow.com/a/35615565
func randString(n int, characterSet string) string {
	if n <= 0 {
		return ""
	}
	result := make([]byte, n)
	bufferSize := int(float64(n) * 1.3)
	for i, j, randomBytes := 0, 0, []byte{}; i < n; j++ {
		if j%bufferSize == 0 {
			randomBytes = secureRandomBytes(bufferSize)
		}
		if idx := int(randomBytes[j%n] & letterIdxMask); idx < len(characterSet) {
			result[i] = characterSet[idx]
			i++
		}
	}

	return string(result)
}

// SecureRandomBytes returns the requested number of bytes using crypto/rand
func secureRandomBytes(length int) []byte {
	var randomBytes = make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal("Unable to generate random bytes")
	}
	return randomBytes
}

// https://stackoverflow.com/a/24894202
func replaceAtIndex(in string, r rune, i int) string {
	out := []rune(in)
	out[i] = r
	return string(out)
}

// AInB check if 'A' exists in 'B'
func AInB(a, b string) bool {
	return a == "" || strings.Contains(
		strings.ToLower(b),
		strings.ToLower(a))
}

// DeduplicateString removes items when there is more than 1 of the same item
// https://stackoverflow.com/a/66751055
func DeduplicateString(s []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range s {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// Sanitize removes some unsafe html
func Sanitize(s string) string {
	return policy.Sanitize(s)
}
