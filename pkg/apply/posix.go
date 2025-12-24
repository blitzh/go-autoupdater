//go:build !windows

package apply

import (
	"context"
	"time"

	"github.com/blitzh/go-autoupdater/pkg/service"
	"github.com/blitzh/go-autoupdater/pkg/util"
)

type PosixApplier struct {
	Retries int
}

func (a PosixApplier) Apply(ctx context.Context, svc service.Controller, currentPath, newPath, oldPath string) (string, error) {
	if a.Retries <= 0 {
		a.Retries = 30
	}

	// stop service/process if provided
	_ = svc.Stop(ctx)

	// best-effort remove old
	_ = util.RemoveWithRetry(oldPath, a.Retries, 200*time.Millisecond)

	// rename current->old (if exists)
	_ = util.RenameWithRetry(currentPath, oldPath, a.Retries, 200*time.Millisecond)

	// rename new->current
	if err := util.RenameWithRetry(newPath, currentPath, a.Retries, 200*time.Millisecond); err != nil {
		// rollback: old->current
		_ = util.RenameWithRetry(oldPath, currentPath, a.Retries, 200*time.Millisecond)
		return "", err
	}

	// start again
	if err := svc.Start(ctx); err != nil {
		// rollback
		_ = svc.Stop(ctx)
		_ = util.RemoveWithRetry(currentPath, a.Retries, 200*time.Millisecond)
		_ = util.RenameWithRetry(oldPath, currentPath, a.Retries, 200*time.Millisecond)
		_ = svc.Start(ctx)
		return "", err
	}

	return oldPath, nil
}
