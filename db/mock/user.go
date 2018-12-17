package mock

import (
	"context"

	core "github.com/ueokande/envoy-playground"
	"github.com/ueokande/envoy-playground/db"
)

func (i *impl) AddUser(ctx context.Context, u core.User) error {
	i.m.Lock()
	defer i.m.Unlock()

	_, ok := i.users[u.Login]
	if ok {
		return db.ErrConflict
	}
	i.users[u.Login] = u
	return nil
}

func (i *impl) GetUser(ctx context.Context, login string) (*core.User, error) {
	i.m.Lock()
	defer i.m.Unlock()

	u, ok := i.users[login]
	if !ok {
		return nil, db.ErrNotFound
	}
	return &u, nil
}

func (i *impl) ListUsers(ctx context.Context) ([]*core.User, error) {
	us := make([]*core.User, len(i.users))
	var n int
	for _, u := range i.users {
		us[n] = &u
		n++
	}
	return us, nil
}

func (i *impl) UpdateUser(ctx context.Context, u core.User) error {
	i.m.Lock()
	defer i.m.Unlock()

	_, ok := i.users[u.Login]
	if !ok {
		return db.ErrNotFound
	}
	i.users[u.Login] = u
	return nil
}

func (i *impl) RemoveUser(ctx context.Context, login string) error {
	i.m.Lock()
	defer i.m.Unlock()

	_, ok := i.users[login]
	if !ok {
		return db.ErrNotFound
	}
	delete(i.users, login)
	return nil
}
