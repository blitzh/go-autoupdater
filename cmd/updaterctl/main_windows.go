//go:build windows

package main

import (
	"github.com/blitzh/go-autoupdater/pkg/apply"
	"github.com/blitzh/go-autoupdater/pkg/service"
)

func main() {
	a := parseArgs()

	// Controller: NSSM by default; allow forcing SC by `--nssm SC`
	var ctrl service.Controller = service.NoopController{}
	if a.svcName != "" {
		if a.nssmPath == "SC" {
			ctrl = service.SCController{ServiceName: a.svcName}
		} else {
			ctrl = service.NSSMController{ServiceName: a.svcName, NSSMPath: a.nssmPath}
		}
	}

	// Applier: Windows needs helper because exe is locked when running
	ap := apply.WindowsHelperApplier{
		ServiceName: a.svcName,
		NSSMPath:    a.nssmPath,
		// HelperPath default = <installDir>\updater-helper.exe (see applier implementation)
	}

	code := runUpdate(a, ctrl, ap)
	// exit code is handled by runtime implicitly
	_ = code
}

func defaultExeName() string { return "agent.exe" }
