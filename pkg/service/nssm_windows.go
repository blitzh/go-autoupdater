//go:build windows

package service

import (
	"context"
	"os/exec"
	"syscall"
)

type NSSMController struct {
	ServiceName string
	NSSMPath    string // optional; if empty will use "nssm" from PATH
}

func (c NSSMController) Stop(ctx context.Context) error {
	return runHideNSSM(ctx, c.bin(), "stop", c.ServiceName)
}
func (c NSSMController) Start(ctx context.Context) error {
	return runHideNSSM(ctx, c.bin(), "start", c.ServiceName)
}
func (c NSSMController) Restart(ctx context.Context) error {
	_ = c.Stop(ctx)
	return c.Start(ctx)
}
func (c NSSMController) String() string { return "nssm:" + c.ServiceName }

func (c NSSMController) bin() string {
	if c.NSSMPath != "" {
		return c.NSSMPath
	}
	return "nssm"
}

func runHideNSSM(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
