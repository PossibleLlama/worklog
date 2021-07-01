package helpers

// Remove removes a single item at index from the slice
// https://stackoverflow.com/a/37335777
func Remove(s []interface{}, i int) []interface{} {
	return append(s[:i], s[i+1:]...)
}
