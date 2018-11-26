// +build linux

package sysctl

import (
	"testing"
)

func TestGet(t *testing.T) {
	got, err := Get("net.ipv4.ip_forward")
	if err != nil {
		t.Fatalf("Could not get sysctl value: %s", err.Error())
	}
	if got != "0" && got != "1" {
		t.Errorf("Expected 0 or 1, got %s", got)
	}
}

func TestGetPattern(t *testing.T) {
	pattern := "^net.ipv4.ipfrag"
	expected := []string{
		"net.ipv4.ipfrag_high_thresh",
		"net.ipv4.ipfrag_low_thresh",
		"net.ipv4.ipfrag_max_dist",
		"net.ipv4.ipfrag_time",
	}
	got, err := GetPattern(pattern)
	if err != nil {
		t.Fatalf("Could not get sysctl values for pattern %s: %s", pattern, err.Error())
	}
	if len(got) < len(expected) {
		// We check if length is < than expected to prevent
		// breaking test cases if new sysctls are added
		t.Fatalf("Expected at least %d matches, got %d. Matches: %+v",
			len(expected), len(got), got)
	}
	for _, k := range expected {
		if _, ok := got[k]; !ok {
			t.Fatalf("Key %s not matched. Matches: %+v", k, got)
		}
	}
}
