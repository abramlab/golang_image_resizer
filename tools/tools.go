//go:build tools

package tools

import (
	_ "github.com/alvaroloes/enumer"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "mvdan.cc/gofumpt"
)
