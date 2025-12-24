//go:build darwin

package main

import (
	"os"

	"github.com/blitzh/go-autoupdater/pkg/apply"
	"github.com/blitzh/go-autoupdater/pkg/service"
)

func main() {
	a := parseArgs()

	var ctrl service.Controller = service.NoopController{}
	if a.launchdLbl != "" {
		ctrl = service.LaunchdController{Label: a.launchdLbl}
	}

	ap := apply.PosixApplier{Retries: 40}

	code := runUpdate(a, ctrl, ap)
	os.Exit(code)
}

func defaultExeName() string { return "agent" }
