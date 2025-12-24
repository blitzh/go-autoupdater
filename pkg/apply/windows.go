//go:build windows

package apply

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"

	"github.com/blitzh/go-autoupdater/pkg/service"
)

type WindowsHelperApplier struct {
	HelperPath  string // optional; if empty, will look in same dir as current exe: updater-helper.exe
	NSSMPath    string // optional; passed to helper
	ServiceName string // must be set if you want stop/start in helper
	Retries     int
}

func (a WindowsHelperApplier) Apply(ctx context.Context, svc service.Controller, currentPath, newPath, oldPath string) (string, error) {
	// Windows: currentPath is locked if service running.
	// We delegate stop/swap/start to helper process.

	helper := a.HelperPath
	if helper == "" {
		// assume helper beside currentPath (install dir)
		helper = filepath.Join(filepath.Dir(currentPath), "updater-helper.exe")
	}
	if _, err := os.Stat(helper); err != nil {
		return "", fmt.Errorf("windows helper not found: %s", helper)
	}

	args := []string{
		"--current", currentPath,
		"--new", newPath,
		"--old", oldPath,
		"--service", a.ServiceName,
	}
	if a.NSSMPath != "" {
		args = append(args, "--nssm", a.NSSMPath)
	}

	cmd := exec.CommandContext(ctx, helper, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	// important: do not inherit stdin blocking service
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Helper will stop/start itself. `svc` here can be noop.
	if err := cmd.Run(); err != nil {
		// give time for filesystem settle
		time.Sleep(500 * time.Millisecond)
		return "", err
	}
	return oldPath, nil
}
