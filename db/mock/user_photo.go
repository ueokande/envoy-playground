package mock

import (
	"context"

	"github.com/ueokande/envoy-playground/db"
)

func (i *impl) GetPhoto(ctx context.Context, login string) (string, error) {
	i.m.Lock()
	defer i.m.Unlock()

	uuid, ok := i.userPhotos[login]
	if !ok {
		return "", db.ErrNotFound
	}
	return uuid, nil
}

func (i *impl) UpdatePhoto(ctx context.Context, login string, uuid string) error {
	i.m.Lock()
	defer i.m.Unlock()

	i.userPhotos[login] = uuid
	return nil
}

func (i *impl) RemovePhoto(ctx context.Context, login string) error {
	i.m.Lock()
	defer i.m.Unlock()

	delete(i.userPhotos, login)
	return nil
}
