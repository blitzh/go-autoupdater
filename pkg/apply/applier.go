package apply

import (
	"context"

	"github.com/blitzh/go-autoupdater/pkg/service"
)

type Applier interface {
	Apply(ctx context.Context, svc service.Controller, currentPath, newPath, oldPath string) (oldBackup string, err error)
}
