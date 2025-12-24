package service

import "context"

type Controller interface {
	Stop(ctx context.Context) error
	Start(ctx context.Context) error
	Restart(ctx context.Context) error
	String() string
}

type NoopController struct{}
