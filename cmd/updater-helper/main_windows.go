//go:build windows

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/blitzh/go-autoupdater/pkg/util"
)

func main() {
	var (
		current = flag.String("current", "", "path to current exe")
		newf    = flag.String("new", "", "path to new exe")
		oldf    = flag.String("old", "", "path to old backup exe")
		svc     = flag.String("service", "", "service name (optional)")
		nssm    = flag.String("nssm", "", "path to nssm.exe (optional)")
	)
	flag.Parse()

	if *current == "" || *newf == "" || *oldf == "" {
		fmt.Println("missing required args: --current --new --old")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	stop := func() error {
		if strings.TrimSpace(*svc) == "" {
			return nil
		}
		if bin := resolveNSSM(*nssm); bin != "" {
			return runHide(ctx, bin, "stop", *svc)
		}
		return runHide(ctx, "sc", "stop", *svc)
	}
	start := func() error {
		if strings.TrimSpace(*svc) == "" {
			return nil
		}
		if bin := resolveNSSM(*nssm); bin != "" {
			return runHide(ctx, bin, "start", *svc)
		}
		return runHide(ctx, "sc", "start", *svc)
	}

	// stop service first
	_ = stop()
	time.Sleep(800 * time.Millisecond)

	// swap
	_ = util.RemoveWithRetry(*oldf, 30, 250*time.Millisecond)
	_ = util.RenameWithRetry(*current, *oldf, 30, 250*time.Millisecond)
	if err := util.RenameWithRetry(*newf, *current, 40, 250*time.Millisecond); err != nil {
		// rollback quickly
		_ = util.RenameWithRetry(*oldf, *current, 40, 250*time.Millisecond)
		fmt.Println("swap failed:", err)
		os.Exit(2)
	}

	// start service
	if err := start(); err != nil {
		// rollback
		_ = stop()
		_ = util.RemoveWithRetry(*current, 30, 250*time.Millisecond)
		_ = util.RenameWithRetry(*oldf, *current, 40, 250*time.Millisecond)
		_ = start()
		fmt.Println("start failed, rolled back:", err)
		os.Exit(3)
	}

	fmt.Println("helper done: update applied")
}

func resolveNSSM(path string) string {
	if strings.TrimSpace(path) != "" {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	if p, err := exec.LookPath("nssm"); err == nil {
		return p
	}
	return ""
}

func runHide(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
