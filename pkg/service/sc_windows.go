//go:build windows

package service

import (
	"context"
	"os/exec"
	"syscall"
)

type SCController struct {
	ServiceName string
}

func (c SCController) Stop(ctx context.Context) error {
	return runHideSC(ctx, "sc", "stop", c.ServiceName)
}
func (c SCController) Start(ctx context.Context) error {
	return runHideSC(ctx, "sc", "start", c.ServiceName)
}
func (c SCController) Restart(ctx context.Context) error { _ = c.Stop(ctx); return c.Start(ctx) }
func (c SCController) String() string                    { return "sc:" + c.ServiceName }

func runHideSC(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
