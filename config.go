package sysctl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const sysctlConfPath = "/etc/sysctl.conf"

// parseConfig reads a sysctl configuration file and parses its content
func parseConfig(path string, out map[string]string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open file: %s", err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parsed := strings.Split(line, "#")[0]
		parsed = strings.Split(parsed, ";")[0]
		parsed = strings.TrimSpace(parsed)
		if parsed == "" {
			continue
		}
		tokens := strings.Split(parsed, " = ")
		if len(tokens) != 2 {
			return fmt.Errorf("could not parse line %s", line)
		}
		k := strings.TrimSpace(tokens[0])
		v := strings.TrimSpace(tokens[1])
		out[k] = v
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %s", err.Error())
	}

	return nil
}

// LoadConfig gets sysctl values from a list of sysctl configuration files.
// The values in the rightmost files take priority.
// If no file is specified, values are read from /etc/sysctl.conf.
func LoadConfig(files ...string) (map[string]string, error) {
	if len(files) == 0 {
		files = []string{sysctlConfPath}
	}
	out := make(map[string]string)
	for _, f := range files {
		if err := parseConfig(f, out); err != nil {
			return nil, fmt.Errorf("could not parse file %s: %s", f, err.Error())
		}
	}
	return out, nil
}

// LoadConfigAndApply sets sysctl values from a list of sysctl configuration files.
// The values in the rightmost files take priority.
// If no file is specified, values are read from /etc/sysctl.conf.
func LoadConfigAndApply(files ...string) error {
	config, err := LoadConfig(files...)
	if err != nil {
		return fmt.Errorf("could not read configuration from files: %s", err.Error())
	}
	for k, v := range config {
		if err := Set(k, v); err != nil {
			return fmt.Errorf("could not set %s = %s: %s", k, v, err.Error())
		}
	}
	return nil
}
