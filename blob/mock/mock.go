package mock

import (
	"bytes"
	"context"
	"io"

	"github.com/ueokande/envoy-playground/blob"
)

type impl struct {
	data map[string][]byte
}

func New() blob.Interface {
	return &impl{data: make(map[string][]byte)}
}

func (i *impl) Get(ctx context.Context, name string) (blob.Object, error) {
	b, ok := i.data[name]
	if !ok {
		return nil, blob.ErrNotFound
	}
	return &object{bytes.NewReader(b)}, nil
}

func (i *impl) Put(ctx context.Context, name string, r io.Reader, size int64) error {
	buf := make([]byte, size)
	_, err := r.Read(buf)
	if err != nil {
		return err
	}
	i.data[name] = buf
	return nil
}

func (i *impl) Delete(ctx context.Context, name string) error {
	delete(i.data, name)
	return nil
}

type object struct {
	*bytes.Reader
}

func (o *object) Close() error {
	return nil
}
