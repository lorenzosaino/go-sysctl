// Package sysctl provides functions wrapping the sysctl interface.
package sysctl

// DefaultPath is the default path to the sysctl virtual files.
const DefaultPath = "/proc/sys/"

var std *Client

func init() {
	std = NewClient(DefaultPath)
}

// Get returns a sysctl from a given key.
func Get(key string) (string, error) {
	return std.Get(key)
}

// GetPattern returns a map of sysctls matching a given pattern
// The pattern uses a POSIX extended regular expression syntax.
// This function matches the same sysctls that the command
// sysctl -a -r <pattern> would return.
func GetPattern(pattern string) (map[string]string, error) {
	return std.GetPattern(pattern)
}

// GetAll returns all sysctls. This is equivalent
// to running the command sysctl -a.
func GetAll() (map[string]string, error) {
	return std.GetAll()
}

// Set updates the value of a sysctl.
func Set(key, value string) error {
	return std.Set(key, value)
}
