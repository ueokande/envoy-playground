package client

import (
	"context"
	"io"

	core "github.com/ueokande/envoy-playground"
)

type Interface interface {
	// ListUsers lists all users
	ListUsers(ctx context.Context) ([]core.User, error)

	// GetUser gets the user
	GetUser(ctx context.Context, login string) (*core.User, error)

	// AddUser adds user
	AddUser(ctx context.Context, login, name string) (*core.User, error)

	// UpdateUser updates the user
	UpdateUser(ctx context.Context, u core.User) (*core.User, error)

	// DeleteUser deletes users
	DeleteUser(ctx context.Context, login string) error

	// UpdateUserPhoto updates user's photo
	UpdateUserPhoto(ctx context.Context, login string, r io.Reader) error

	// GetUserPhoto gets user's photo
	GetUserPhoto(ctx context.Context, w io.Writer, login string) error

	// DeleteUserPhoto deletes user's photo
	DeleteUserPhoto(ctx context.Context, login string) error
}
