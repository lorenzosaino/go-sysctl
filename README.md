# Go Sysctl

Golang wrapper around the sysctl interface.

## Documentation

See [Go doc](https://godoc.org/github.com/lorenzosaino/go-sysctl).

## Example

```go
import sysctl "github.com/lorenzosaino/go-sysctl"

// Get value of a single sysctl
// This is equivalent to running "sysctl <key>"
val, err := sysctl.Get("net.ipv4.ip_forward")

// Get the values of all sysctls matching a given pattern
// This is equivalent to running "sysctl -a -r <pattern>"
vals, err := sysctl.GetPattern("net.ipv4.ipfrag")

// Get the values of all sysctls
// This is equivalent to running "sysctl -a"
vals, err = sysctl.GetAll()

// Set the value of a sysctl
// This is equivalent to running "sysctl -w <key>=<value>"
err = sysctl.Set("net.ipv4.ip_forward", "1")

// Set sysctl values from configuration file
// This is equivalent to running "sysctl -p <config-file>"
err = LoadConfigAndApply("/etc/sysctl.conf")
```
