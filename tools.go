//go:build tools
// +build tools

package sysctl

import (
	_ "golang.org/x/lint/golint"
	_ "golang.org/x/vuln/cmd/govulncheck"
	_ "honnef.co/go/tools/cmd/staticcheck"
)
