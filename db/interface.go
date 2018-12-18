package db

import (
	"context"
	"errors"

	core "github.com/ueokande/envoy-playground"
)

type Interface interface {
	// AddUser adds an user
	AddUser(ctx context.Context, u core.User) error
	// GetUser returns the user of login
	GetUser(ctx context.Context, login string) (*core.User, error)
	// ListUsers returns all users
	ListUsers(ctx context.Context) ([]*core.User, error)
	// UpdateUser updates a user
	UpdateUser(ctx context.Context, u core.User) error
	// RemoveUser removes a user of the login
	RemoveUser(ctx context.Context, login string) error

	// GetPhoto returns user's photo
	GetPhoto(ctx context.Context, login string) (string, error)
	// UpdatePhoto updates login's photo
	UpdatePhoto(ctx context.Context, login string, uuid string) error
	// RemovePhoto removes user's photo
	RemovePhoto(ctx context.Context, login string) error

	Close() error
}

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)
