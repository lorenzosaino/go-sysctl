// Package sysctl provides functions wrapping the sysctl interface.
package sysctl

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func isFileReadable(info os.FileInfo) bool {
	// other users have read permissions
	// this is not completely accurate because
	// we should also check if the UID or GID of
	// the file match those of the current user
	// and if group or user have read permissions
	return info.Mode()&(1<<2) != 0
}

func readFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func writeFile(path, value string) error {
	return os.WriteFile(path, []byte(value), 0644)
}

// Client is a client for reading and writing sysctls
type Client struct {
	path string
}

// NewClient returns a new Client.
// The path argument is the base path containing all sysctl virtual files.
// By default this is DefaultPath, but there may be cases where you may want
// to use a different path, e.g. for tests or if procfs path is mounted
// to a different path.
func NewClient(path string) *Client {
	return &Client{path: path}
}

func (c *Client) pathFromKey(key string) string {
	return filepath.Join(c.path, strings.Replace(key, ".", "/", -1))
}

func (c *Client) keyFromPath(path string) string {
	subPath := strings.TrimPrefix(path, c.path)
	return strings.Replace(subPath, "/", ".", -1)
}

// Get returns a sysctl from a given key.
func (c *Client) Get(key string) (string, error) {
	return readFile(c.pathFromKey(key))
}

// GetPattern returns a map of sysctls matching a given pattern
// The pattern uses a POSIX extended regular expression syntax.
// This function matches the same sysctls that the command
// sysctl -a -r <pattern> would return.
func (c *Client) GetPattern(pattern string) (map[string]string, error) {
	re, err := regexp.CompilePOSIX(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %v", err)
	}
	res := make(map[string]string)
	err = filepath.Walk(c.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing sysctl path: %v", err)
		}
		if info.IsDir() {
			return nil
		}
		key := c.keyFromPath(path)
		if !re.MatchString(key) {
			return nil
		}
		if !isFileReadable(info) {
			return nil
		}
		val, err := readFile(path)
		if err != nil {
			return fmt.Errorf("error reading %s: %v", path, err)
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
func (c *Client) GetAll() (map[string]string, error) {
	return c.GetPattern("")
}

// Set updates the value of a sysctl.
func (c *Client) Set(key, value string) error {
	return writeFile(c.pathFromKey(key), value)
}
