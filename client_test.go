package sysctl

import (
	"testing"
)

func TestClient_pathFromKey(t *testing.T) {
	cases := []struct {
		base     string
		in       string
		expected string
	}{
		{
			base:     "/proc/sys/",
			in:       "net.ipv4.ip_forward",
			expected: "/proc/sys/net/ipv4/ip_forward",
		},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			cl := NewClient(c.base)
			got := cl.pathFromKey(c.in)
			if got != c.expected {
				t.Fatalf("Expected: %s. Got: %s", c.expected, got)
			}
		})
	}
}

func TestClient_keyFromPath(t *testing.T) {
	cases := []struct {
		base     string
		in       string
		expected string
	}{
		{
			base:     "/proc/sys/",
			in:       "/proc/sys/net/ipv4/ip_forward",
			expected: "net.ipv4.ip_forward",
		},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			cl := NewClient(c.base)
			got := cl.keyFromPath(c.in)
			if got != c.expected {
				t.Fatalf("Expected: %s. Got: %s", c.expected, got)
			}
		})
	}
}
