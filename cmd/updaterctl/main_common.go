package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/blitzh/go-autoupdater/pkg/apply"
	"github.com/blitzh/go-autoupdater/pkg/service"
	"github.com/blitzh/go-autoupdater/pkg/source"
	"github.com/blitzh/go-autoupdater/pkg/updater"
	"github.com/blitzh/go-autoupdater/pkg/util"
)

type cliArgs struct {
	manifestURL string
	installDir  string
	exeName     string
	curVer      string
	logFile     string

	// service-specific (may be unused on some OS builds)
	svcName     string
	nssmPath    string
	systemdUnit string
	launchdLbl  string

	timeout time.Duration
}

func parseArgs() cliArgs {
	var a cliArgs

	flag.StringVar(&a.manifestURL, "manifest", "", "manifest.json url")
	flag.StringVar(&a.installDir, "dir", ".", "install directory")
	flag.StringVar(&a.exeName, "exe", "", "executable name (e.g. agent.exe / agent)")
	flag.StringVar(&a.curVer, "current", "", "current version (optional)")
	flag.StringVar(&a.logFile, "log", "", "log file path (optional)")

	flag.StringVar(&a.svcName, "service", "", "service name (windows) (optional)")
	flag.StringVar(&a.nssmPath, "nssm", "", "path to nssm.exe (windows) (optional; use \"SC\" to force sc.exe)")
	flag.StringVar(&a.systemdUnit, "systemd", "", "systemd unit (linux) (optional)")
	flag.StringVar(&a.launchdLbl, "launchd", "", "launchd label (darwin) (optional)")

	flag.DurationVar(&a.timeout, "timeout", 120*time.Second, "update timeout")

	flag.Parse()
	return a
}

func runUpdate(a cliArgs, ctrl service.Controller, ap apply.Applier) int {
	if a.manifestURL == "" {
		fmt.Println("missing --manifest")
		return 2
	}

	if a.exeName == "" {
		a.exeName = defaultExeName()
	}

	if a.logFile == "" {
		a.logFile = filepath.Join(a.installDir, "updaterctl.log")
	}
	logger := util.NewLogger(a.logFile)

	src := source.NewHTTPManifestSource(a.manifestURL)

	u := updater.New(updater.Config{
		CurrentVersion: a.curVer,
		InstallDir:     a.installDir,
		ExeName:        a.exeName,
		Source:         src,
		Service:        ctrl,
		Applier:        ap,
		Logger:         logger,
	})

	ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
	defer cancel()

	res, err := u.Update(ctx)
	if err != nil {
		logger.Printf("update failed: %v", err)
		fmt.Println("update failed:", err)
		return 1
	}

	if !res.DidUpdate {
		logger.Printf("no update. remote=%s", res.RemoteVersion)
		fmt.Println("no update. remote=", res.RemoteVersion)
		return 0
	}

	logger.Printf("updated OK -> %s, backup=%s", res.RemoteVersion, res.OldBackupPath)
	fmt.Println("updated OK ->", res.RemoteVersion)
	return 0
}
