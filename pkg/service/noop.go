package service

import (
	"context"
)

func (n NoopController) Stop(ctx context.Context) error    { return nil }
func (n NoopController) Start(ctx context.Context) error   { return nil }
func (n NoopController) Restart(ctx context.Context) error { return nil }
func (n NoopController) String() string                    { return "noop" }
