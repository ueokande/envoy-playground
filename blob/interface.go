package blob

import (
	"context"
	"errors"
	"io"
)

type Object interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type Interface interface {
	Get(ctx context.Context, name string) (Object, error)

	Put(ctx context.Context, name string, r io.Reader) error

	Delete(ctx context.Context, name string) error
}

var ErrNotFound = errors.New("not found")
