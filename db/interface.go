package db

import (
	"context"

	core "github.com/ueokande/envoy-playground"
)

type Interface interface {
	AddUser(ctx context.Context, u core.User) error

	GetUser(ctx context.Context, login string) (*core.User, error)

	ListUsers(ctx context.Context) ([]*core.User, error)

	UpdateUser(ctx context.Context, u core.User) error

	RemoveUser(ctx context.Context, login string) error

	Close() error
}
