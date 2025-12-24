//go:build linux

package service

import (
	"context"
	"os/exec"
)

type SystemdController struct {
	Unit string // e.g. "agent.service"
}

func (s SystemdController) Stop(ctx context.Context) error {
	return exec.CommandContext(ctx, "systemctl", "stop", s.Unit).Run()
}
func (s SystemdController) Start(ctx context.Context) error {
	return exec.CommandContext(ctx, "systemctl", "start", s.Unit).Run()
}
func (s SystemdController) Restart(ctx context.Context) error {
	return exec.CommandContext(ctx, "systemctl", "restart", s.Unit).Run()
}
func (s SystemdController) String() string { return "systemd:" + s.Unit }
