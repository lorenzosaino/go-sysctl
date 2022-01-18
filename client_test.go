package sysctl

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewClient(t *testing.T) {
	cases := []struct {
		name string
		path string
		ok   bool
	}{
		{
			name: "dir not found",
			path: "testdata/client/not-found",
			ok:   false,
		},
		{
			name: "file trailing slash",
			path: "testdata/client/ok/f/",
			ok:   false,
		},
		{
			name: "not a dir",
			path: "testdata/client/ok/f",
			ok:   false,
		},
		{
			name: "ok",
			path: "testdata/client/ok",
			ok:   true,
		},
		{
			name: "ok trailing slash",
			path: "testdata/client/ok/",
			ok:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl, err := NewClient(c.path)
			if c.ok && err != nil {
				t.Fatalf("error parsing: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("expected error but it succeeded")
			}
			if err != nil {
				t.Logf("err: %v", err)
				return
			}
			if cl == nil {
				t.Fatal("client unexpectedly nil")
			}
		})
	}
}

func TestClient_pathFromKey(t *testing.T) {
	cases := []struct {
		base     string
		in       string
		expected string
	}{
		{
			base:     "/a/b/",
			in:       "net.ipv4.ip_forward",
			expected: "/a/b/net/ipv4/ip_forward",
		},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			cl := Client{path: c.base}
			got := cl.pathFromKey(c.in)
			if got != c.expected {
				t.Fatalf("expected: %s. Got: %s", c.expected, got)
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
			base:     "/a/b/",
			in:       "/a/b/net/ipv4/ip_forward",
			expected: "net.ipv4.ip_forward",
		},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			cl := Client{path: c.base}
			got := cl.keyFromPath(c.in)
			if got != c.expected {
				t.Fatalf("expected: %s. Got: %s", c.expected, got)
			}
		})
	}
}

func TestClientGet(t *testing.T) {
	cases := []struct {
		name string
		path string
		key  string
		val  string
		ok   bool
	}{
		{
			name: "missing key",
			path: "testdata/client/nok",
			key:  "missing",
			ok:   false,
		},
		{
			name: "ok root file",
			path: "testdata/client/ok",
			key:  "f",
			val:  "value of f",
			ok:   true,
		},
		{
			name: "ok one level",
			path: "testdata/client/ok",
			key:  "d.f",
			val:  "value of d.f",
			ok:   true,
		},
		{
			name: "ok two levels",
			path: "testdata/client/ok",
			key:  "d.d.f1",
			val:  "value of d.d.f1",
			ok:   true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cl, err := NewClient(c.path)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			got, err := cl.Get(c.key)
			if c.ok && err != nil {
				t.Fatalf("error parsing: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("expected error but it succeeded")
			}
			if err != nil {
				t.Logf("err: %v", err)
				return
			}
			if got != c.val {
				t.Fatalf("expected: %s. Got: %s", c.val, got)
			}
		})
	}
}

func TestClientGetPattern(t *testing.T) {
	cases := []struct {
		name    string
		path    string
		pattern string
		res     map[string]string
		ok      bool
	}{
		{
			name: "empty",
			res:  map[string]string{},
			ok:   true,
		},
		{
			name:    "invalid pattern",
			pattern: "[[",
		},
		{
			name:    "match all",
			path:    "testdata/client/ok/",
			pattern: "",
			res: map[string]string{
				"f":      "value of f",
				"d.f":    "value of d.f",
				"d.d.f1": "value of d.d.f1",
				"d.d.f2": "value of d.d.f2",
			},
			ok: true,
		},
		{
			name:    "ok star match",
			path:    "testdata/client/ok/",
			pattern: "d.*",
			res: map[string]string{
				"d.f":    "value of d.f",
				"d.d.f1": "value of d.d.f1",
				"d.d.f2": "value of d.d.f2",
			},
			ok: true,
		},
		{
			name:    "ok single char match",
			path:    "testdata/client/ok/",
			pattern: "d.d.f?",
			res: map[string]string{
				"d.d.f1": "value of d.d.f1",
				"d.d.f2": "value of d.d.f2",
			},
			ok: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := c.path
			if path == "" {
				path = t.TempDir()
				defer os.RemoveAll(path)
			}
			cl, err := NewClient(path)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			got, err := cl.GetPattern(c.pattern)
			if c.ok && err != nil {
				t.Fatalf("error parsing: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("expected error but it succeeded")
			}
			if err != nil {
				t.Logf("err: %v", err)
				return
			}
			if diff := cmp.Diff(c.res, got); diff != "" {
				t.Fatalf("unexpected output (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClientGetAll(t *testing.T) {
	cases := []struct {
		name string
		path string
		res  map[string]string
		ok   bool
	}{
		{
			name: "empty",
			res:  map[string]string{},
			ok:   true,
		},
		{
			name: "testdata/client/ok/",
			path: "testdata/client/ok/",
			res: map[string]string{
				"f":      "value of f",
				"d.f":    "value of d.f",
				"d.d.f1": "value of d.d.f1",
				"d.d.f2": "value of d.d.f2",
			},
			ok: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := c.path
			if path == "" {
				path = t.TempDir()
				defer os.RemoveAll(path)
			}
			cl, err := NewClient(path)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			got, err := cl.GetAll()
			if c.ok && err != nil {
				t.Fatalf("error parsing: %v", err)
			}
			if !c.ok && err == nil {
				t.Fatal("expected error but it succeeded")
			}
			if err != nil {
				t.Logf("err: %v", err)
				return
			}
			if diff := cmp.Diff(c.res, got); diff != "" {
				t.Fatalf("unexpected output (-want +got):\n%s", diff)
			}
		})
	}
}

func TestClientSet(t *testing.T) {
	cases := []struct {
		name  string
		files []string
		keys  map[string]string
		ok    bool
	}{
		{
			name: "empty",
			ok:   true,
		},
		{
			name: "ok",
			files: []string{
				"a",
				"b/a",
				"b/b/a",
			},
			keys: map[string]string{
				"a":     "value of a",
				"b.a":   "value of b.a",
				"b.b.a": "value of b.b.a",
			},
			ok: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			path := t.TempDir()
			defer os.RemoveAll(path)
			createTestFiles(t, path, c.files)

			cl, err := NewClient(path)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}
			for k, v := range c.keys {
				err := cl.Set(k, v)
				if c.ok && err != nil {
					t.Fatalf("could not set %s=%s: %v", k, v, err)
				}
				if !c.ok && err == nil {
					t.Fatal("expected error but it succeeded")
				}
				if err != nil {
					t.Logf("err: %v", err)
					return
				}
			}
			for k, v := range c.keys {
				got, err := cl.Get(k)
				if err != nil {
					t.Fatalf("could not get key %s: %v", k, err)
				}
				if got != v {
					t.Fatalf("got wrong value for key %s: expected: %s, got %s", k, v, got)
				}
			}
		})
	}
}

func TestClientLoadConfigAndApply(t *testing.T) {
	cases := []struct {
		name     string
		files    []string
		config   string
		expected map[string]string
		ok       bool
	}{
		{
			name:   "empty",
			config: "testdata/client/config-empty.conf",
			ok:     true,
		},
		{
			name:   "ok",
			config: "testdata/client/config-ok.conf",
			files: []string{
				"a",
				"b/a",
				"b/b/a",
			},
			expected: map[string]string{
				"a":     "value of a",
				"b.a":   "value of b.a",
				"b.b.a": "value of b.b.a",
			},
			ok: true,
		},
		{
			name:   "missing config file",
			config: "testdata/client/missing-config-file.conf",
		},
		{
			name:   "missing keys",
			config: "testdata/client/config-missing-keys.conf",
			files: []string{
				"a",
				"b/a",
				"b/b/a",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {

			path := t.TempDir()
			defer os.RemoveAll(path)
			createTestFiles(t, path, c.files)

			cl, err := NewClient(path)
			if err != nil {
				t.Fatalf("could not create client: %v", err)
			}

			err = cl.LoadConfigAndApply(c.config)
			if c.ok && err != nil {
				t.Fatalf("could not load and apply config: %s", err)
			}
			if !c.ok && err == nil {
				t.Fatal("expected error but it succeeded")
			}
			if err != nil {
				t.Logf("err: %v", err)
				return
			}

			for k, v := range c.expected {
				got, err := cl.Get(k)
				if err != nil {
					t.Fatalf("could not get key %s: %v", k, err)
				}
				if got != v {
					t.Fatalf("got wrong value for key %s: expected: %s, got %s", k, v, got)
				}
			}
		})
	}
}

func createTestFiles(t *testing.T, base string, paths []string) {
	t.Helper()
	for _, p := range paths {
		p := filepath.Join(base, p)
		dir := filepath.Dir(p)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			t.Fatalf("could not create dir %s: %v", dir, err)
		}
		f, err := os.Create(p)
		if err != nil {
			t.Fatalf("could not create file %s: %v", p, err)
		}
		_ = f.Close()
	}
}
