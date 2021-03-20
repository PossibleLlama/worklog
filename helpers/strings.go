package helpers

import (
	"math/rand"
	"time"
	"unsafe"
)

var src = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	alphabeticBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hexadecimalBytes = "abcdef1234567890"
	letterIdxBits    = 6
	letterIdxMask    = 1<<letterIdxBits - 1
	letterIdxMax     = 63 / letterIdxBits
)

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

// From https://stackoverflow.com/a/31832326
func randString(n int, characterSet string) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(characterSet) {
			b[i] = characterSet[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
