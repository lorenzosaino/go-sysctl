// Package sysctl provides functions wrapping the sysctl interface.
package sysctl

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const sysctlBase = "/proc/sys/"

func pathFromKey(key string) string {
	return filepath.Join(sysctlBase, strings.Replace(key, ".", "/", -1))
}

func keyFromPath(path string) string {
	subPath := strings.TrimPrefix(path, sysctlBase)
	return strings.Replace(subPath, "/", ".", -1)
}

func isFileReadable(info os.FileInfo) bool {
	// other users have read permissions
	// this is not completely accurate because
	// we should also check if the UID or GID of
	// the file match those of the current user
	// and if group or user have read permissions
	return info.Mode()&(1<<2) != 0
}

func readFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func writeFile(path, value string) error {
	return ioutil.WriteFile(path, []byte(value), 0644)
}

// Get returns a sysctl from a given key.
func Get(key string) (string, error) {
	return readFile(pathFromKey(key))
}

// GetPattern returns a map of sysctls matching a given pattern
// The pattern uses a POSIX extended regular expression syntax.
// This function matches the same sysctls that the command
// sysctl -a -r <pattern> would return.
func GetPattern(pattern string) (map[string]string, error) {
	re, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %s", err.Error())
	}
	res := make(map[string]string)
	err = filepath.Walk(sysctlBase, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing sysctl path: %s", err.Error())
		}
		if info.IsDir() {
			return nil
		}
		key := keyFromPath(path)
		if !re.MatchString(key) {
			return nil
		}
		if !isFileReadable(info) {
			return nil
		}
		val, err := readFile(path)
		if err != nil {
			return fmt.Errorf("error reading %s: %s", path, err.Error())
		}
		res[key] = val
		return nil
	})
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetAll returns all sysctls. This is equivalent
// to running the command sysctl -a.
func GetAll() (map[string]string, error) {
	return GetPattern("")
}

// Set updates the value of a sysctl.
func Set(key, value string) error {
	return writeFile(pathFromKey(key), value)
}
