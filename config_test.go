package sysctl

import (
	"reflect"
	"testing"
)

func TestParseConfig(t *testing.T) {
	cases := []struct {
		path string
		err  bool
		out  map[string]string
	}{
		{
			path: "testdata/sysctl-correct.conf",
			err:  false,
			out: map[string]string{
				"kernel.domainname": "example.com",
				"kernel.modprobe":   "/sbin/mod probe",
				"kernel.hostname":   "example.com",
			},
		},
		{
			path: "testdata/sysctl-empty.conf",
			err:  false,
			out:  map[string]string{},
		},
		{
			path: "testdata/sysctl-error.conf",
			err:  true,
		},
	}
	for _, c := range cases {
		out := make(map[string]string)
		if err := parseConfig(c.path, out); err != nil {
			if !c.err {
				t.Errorf("expected error when parsing %s but it succeeded", c.path)
			}
			continue
		}
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("unexpected output from %s. Expected: %s, got: %s", c.path, c.out, out)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	cases := []struct {
		paths []string
		err   bool
		out   map[string]string
	}{
		{
			paths: []string{
				"testdata/sysctl-a.conf",
				"testdata/sysctl-b.conf",
			},
			err: false,
			out: map[string]string{
				"kernel.domainname": "b.com",
			},
		},
	}
	for _, c := range cases {
		out, err := LoadConfig(c.paths...)
		if err != nil {
			if !c.err {
				t.Errorf("expected error when parsing %s but it succeeded", c.paths)
			}
			continue
		}
		if !reflect.DeepEqual(out, c.out) {
			t.Errorf("unexpected output from %s. Expected: %s, got: %s", c.paths, c.out, out)
		}
	}
}
