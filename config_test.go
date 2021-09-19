package sysctl

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_parseConfig(t *testing.T) {
	cases := []struct {
		name string
		path string
		ok   bool
		out  map[string]string
	}{
		{
			name: "ok",
			path: "testdata/sysctl-correct.conf",
			ok:   true,
			out: map[string]string{
				"kernel.domainname": "example.com",
				"kernel.modprobe":   "/sbin/mod probe",
				"kernel.hostname":   "example.com",
			},
		},
		{
			name: "empty",
			path: "testdata/sysctl-empty.conf",
			ok:   true,
			out:  map[string]string{},
		},
		{
			name: "only-comments",
			path: "testdata/sysctl-only-comments.conf",
			ok:   true,
			out:  map[string]string{},
		},
		{
			name: "malformatted",
			path: "testdata/sysctl-error.conf",
			ok:   false,
		},
		{
			name: "not-found",
			path: "testdata/not-found",
			ok:   false,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out := make(map[string]string)
			err := parseConfig(c.path, out)
			if c.ok && err != nil {
				t.Fatalf("error parsing: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatalf("expected error when parsing %s but it succeeded", c.path)
			}
			if err != nil {
				t.Logf("err: %v", err)
				return
			}
			if diff := cmp.Diff(c.out, out); diff != "" {
				t.Fatalf("unexpected output from %s (-want +got):\n%s", c.path, diff)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	cases := []struct {
		name  string
		paths []string
		ok    bool
		out   map[string]string
	}{
		{
			name: "empty",
			ok:   false,
		},
		{
			name:  "not-found",
			paths: []string{"testdata/not-found"},
			ok:    false,
		},
		{
			name: "ok",
			paths: []string{
				"testdata/sysctl-a.conf",
				"testdata/sysctl-b.conf",
			},
			ok: true,
			out: map[string]string{
				"kernel.domainname": "b.com",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := LoadConfig(c.paths...)
			if c.ok && err != nil {
				t.Fatalf("error parsing: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatalf("expected error when parsing [%s] but it succeeded", strings.Join(c.paths, ", "))
			}
			if err != nil {
				t.Logf("err: %v", err)
				return
			}
			if diff := cmp.Diff(c.out, out); diff != "" {
				t.Fatalf("unexpected output from %s (-want +got):\n%s", c.paths, diff)
			}
		})
	}
}
