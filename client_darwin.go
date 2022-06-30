package sysctl

import (
	"fmt"
	"io"
)

// DefaultPath is the default path to the sysctl virtual files.
const DefaultPath = "/proc/sys/"

var std *Client

func init() {
	std = &Client{path: DefaultPath}
}

type execFn func(string) (io.Reader, error)

// Client is a client for reading and writing sysctls
type Client struct {
	path string
	f    execFn
}

// NewClient returns a new Client.
// The path argument is the base path containing all sysctl virtual files.
// By default this is DefaultPath, but there may be cases where you may want
// to use a different path, e.g. for tests or if procfs path is mounted
// to a different path.
func NewClient(path string) (*Client, error) {
	return &Client{}, nil
}

// Get returns a sysctl from a given key.
func (c *Client) Get(key string) (string, error) {
	// sysctl -e
	return "", nil
}

// GetPattern returns a map of sysctls matching a given pattern
// The pattern uses a POSIX extended regular expression syntax.
// This function matches the same sysctls that the command
// sysctl -a -r <pattern> would return.
func (c *Client) GetPattern(pattern string) (map[string]string, error) {
	return nil, nil
}

// GetAll returns all sysctls. This is equivalent
// to running the command sysctl -a.
func (c *Client) GetAll() (map[string]string, error) {
	// sysctl -a -e
	return nil, nil
}

// Set updates the value of a sysctl.
func (c *Client) Set(key, value string) error {
	return nil
}

// LoadConfigAndApply sets sysctl values from a list of sysctl configuration files.
// The values in the rightmost files take priority.
// If no file is specified, values are read from /etc/sysctl.conf.
func (c *Client) LoadConfigAndApply(files ...string) error {
	config, err := LoadConfig(files...)
	if err != nil {
		return fmt.Errorf("could not read configuration from files: %v", err)
	}
	for k, v := range config {
		if err := c.Set(k, v); err != nil {
			return fmt.Errorf("could not set %s = %s: %v", k, v, err)
		}
	}
	return nil
}
