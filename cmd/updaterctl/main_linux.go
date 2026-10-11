//go:build linux

package main

import (
	"os"

	"github.com/blitzh/go-autoupdater/pkg/apply"
	"github.com/blitzh/go-autoupdater/pkg/service"
)

func main() {
	a := parseArgs()

	var ctrl service.Controller = service.NoopController{}
	if a.systemdUnit != "" {
		ctrl = service.SystemdController{Unit: a.systemdUnit}
	}

	ap := apply.PosixApplier{Retries: 40}

	code := runUpdate(a, ctrl, ap)
	os.Exit(code)
}

func defaultExeName() string { return "agent" }
