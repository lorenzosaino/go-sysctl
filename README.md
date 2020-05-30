# Go Sysctl

[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/lorenzosaino/go-sysctl)
[![Build](https://github.com/lorenzosaino/go-sysctl/workflows/Build/badge.svg)](https://github.com/lorenzosaino/go-sysctl/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/lorenzosaino/go-sysctl)](https://goreportcard.com/report/github.com/lorenzosaino/go-sysctl)
[![License](https://img.shields.io/github/license/lorenzosaino/go-sysctl.svg)](https://github.com/lorenzosaino/go-sysctl/blob/master/LICENSE)

Go wrapper around the sysctl interface.

## Documentation

See [Go doc](https://pkg.go.dev/github.com/lorenzosaino/go-sysctl?tab=doc).

## Usage

```go
import sysctl "github.com/lorenzosaino/go-sysctl"

var (
    val string
    vals map[string]string
    err error
)

// Get value of a single sysctl
// This is equivalent to running "sysctl <key>"
val, err = sysctl.Get("net.ipv4.ip_forward")

// Get the values of all sysctls matching a given pattern
// This is equivalent to running "sysctl -a -r <pattern>"
vals, err = sysctl.GetPattern("net.ipv4.ipfrag")

// Get the values of all sysctls
// This is equivalent to running "sysctl -a"
vals, err = sysctl.GetAll()

// Set the value of a sysctl
// This is equivalent to running "sysctl -w <key>=<value>"
err = sysctl.Set("net.ipv4.ip_forward", "1")

// Set sysctl values from configuration file
// This is equivalent to running "sysctl -p <config-file>"
err = sysctl.LoadConfigAndApply("/etc/sysctl.conf")
```

## License

[BSD 3-clause](LICENSE)
