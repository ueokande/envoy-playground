package mysql

import (
	"context"

	core "github.com/ueokande/envoy-playground"
)

func (i *impl) AddUser(ctx context.Context, u core.User) error {
	_, err := i.db.Exec(`INSERT INTO user(login, name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())`, u.Login, u.Name)
	return err
}

func (i *impl) GetUser(ctx context.Context, login string) (*core.User, error) {
	var u core.User
	err := i.db.QueryRow(`SELECT id,login,name,created_at,updated_at FROM user WHERE login=?`, login).
		Scan(&(u.ID), &(u.Login), &(u.Name), &(u.CreatedAt), &(u.UpdatedAt))
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (i *impl) ListUsers(ctx context.Context) ([]*core.User, error) {
	rows, err := i.db.Query(`SELECT id,login,name,created_at,updated_at FROM user`)
	if err != nil {
		return nil, err
	}

	var us []*core.User
	for rows.Next() {
		var u core.User
		err := rows.Scan(&(u.ID), &(u.Login), &(u.Name), &(u.CreatedAt), &(u.UpdatedAt))
		if err != nil {
			return nil, err
		}
		us = append(us, &u)
	}
	return us, nil
}

func (i *impl) UpdateUser(ctx context.Context, u core.User) error {
	_, err := i.db.Exec(`UPDATE user SET login=?, name=?, updated_at=NOW() WHERE login=?`, u.Login, u.Name, u.Login)
	return err
}

func (i *impl) RemoveUser(ctx context.Context, login string) error {
	_, err := i.db.Exec(`DELETE FROM user WHERE Login=?`, login)
	return err
}
