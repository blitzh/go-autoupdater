//go:build darwin

package service

import (
	"context"
	"os/exec"
)

type LaunchdController struct {
	Label string // e.g. "com.your.agent"
}

func (l LaunchdController) Stop(ctx context.Context) error {
	// best-effort; many setups use bootstrap/bootout; keep simple
	_ = exec.CommandContext(ctx, "launchctl", "stop", l.Label).Run()
	return nil
}
func (l LaunchdController) Start(ctx context.Context) error {
	_ = exec.CommandContext(ctx, "launchctl", "start", l.Label).Run()
	return nil
}
func (l LaunchdController) Restart(ctx context.Context) error {
	_ = l.Stop(ctx)
	return l.Start(ctx)
}
func (l LaunchdController) String() string { return "launchd:" + l.Label }
