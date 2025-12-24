//go:build !windows && !linux && !darwin

package main

import (
	"os"

	"github.com/blitzh/go-autoupdater/pkg/apply"
	"github.com/blitzh/go-autoupdater/pkg/service"
)

func main() {
	a := parseArgs()
	ctrl := service.NoopController{}
	ap := apply.PosixApplier{Retries: 40}
	code := runUpdate(a, ctrl, ap)
	os.Exit(code)
}

func defaultExeName() string { return "agent" }
