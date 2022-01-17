//go:build linux
// +build linux

package sysctl

import (
	"os/user"
	"testing"
)

func TestGet(t *testing.T) {
	cases := []struct {
		name string
		skip bool
	}{
		{
			name: "fs.protected_fifos",
			skip: !isUserRoot(),
		},
		{
			name: "net.ipv4.ip_forward",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.skip {
				t.Skip("skipping test")
			}

			got, err := Get(c.name)
			if err != nil {
				t.Fatalf("Could not get sysctl value: %s", err.Error())
			}
			if got != "0" && got != "1" {
				t.Fatalf("expected 0 or 1, got %s", got)
			}
		})
	}
}

func TestGetPattern(t *testing.T) {
	cases := []struct {
		pattern string
		matches []string
		skip    bool
	}{
		{
			pattern: "^fs.protected_",
			matches: []string{
				"fs.protected_fifos",
				"fs.protected_hardlinks",
				"fs.protected_regular",
				"fs.protected_symlinks",
			},
			skip: !isUserRoot(),
		},
		{
			pattern: "^net.ipv4.ipfrag",
			matches: []string{
				"net.ipv4.ipfrag_high_thresh",
				"net.ipv4.ipfrag_low_thresh",
				"net.ipv4.ipfrag_max_dist",
				"net.ipv4.ipfrag_time",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.pattern, func(t *testing.T) {
			if c.skip {
				t.Skip("skipping test")
			}

			got, err := GetPattern(c.pattern)
			if err != nil {
				t.Fatalf("could not get sysctl values for pattern %s: %v", c.pattern, err)
			}
			if len(got) < len(c.matches) {
				// We check if length is < than expected to prevent
				// breaking test cases if new sysctls are added
				t.Fatalf("expected at least %d matches, got %d. Matches: %+v",
					len(c.matches), len(got), got)
			}
			for _, k := range c.matches {
				if _, ok := got[k]; !ok {
					t.Fatalf("key %s not matched. Matches: %+v", k, got)
				}
			}
		})
	}
}

func isUserRoot() bool {
	u, err := user.Current()
	if err != nil {
		return false
	}
	return u.Username == "root"
}
