package mysql

import (
	"context"

	"github.com/ueokande/envoy-playground/db"
)

func (i *impl) GetPhoto(ctx context.Context, login string) (string, error) {
	var uuid string
	err := i.db.QueryRow(`SELECT user_photo.uuid from user_photo INNER JOIN user on user_photo.user_id = user.id WHERE user.login=?`, login).
		Scan(&uuid)

	if isNotFound(err) {
		return "", db.ErrNotFound
	} else if err != nil {
		return "", err
	}
	return uuid, nil
}

func (i *impl) UpdatePhoto(ctx context.Context, login string, uuid string) error {
	res, err := i.db.Exec(`REPLACE INTO user_photo(user_id,uuid) SELECT user.id,? FROM user WHERE user.login=?;`, uuid, login)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return db.ErrNotFound
	}
	return err
}

func (i *impl) RemovePhoto(ctx context.Context, login string) error {
	res, err := i.db.Exec(`DELETE user_photo FROM user_photo INNER JOIN user ON user_photo.user_id = user.id WHERE user.login=?`, login)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return db.ErrNotFound
	}
	return nil
}
