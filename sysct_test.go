package sysctl

import (
	"testing"
)

func TestPathFromKey(t *testing.T) {
	in := "net.ipv4.ip_forward"
	expected := "/proc/sys/net/ipv4/ip_forward"
	got := pathFromKey(in)
	if got != expected {
		t.Fatalf("Expected: %s. Got: %s", expected, got)
	}
}

func TestKeyFromPath(t *testing.T) {
	in := "/proc/sys/net/ipv4/ip_forward"
	expected := "net.ipv4.ip_forward"
	got := keyFromPath(in)
	if got != expected {
		t.Fatalf("Expected: %s. Got: %s", expected, got)
	}
}
