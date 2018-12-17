package core

import (
	"errors"
	"regexp"
	"time"
)

var loginRegex = regexp.MustCompile(`^[a-z_]{5}[a-z_]*$`)
var nameRegex = regexp.MustCompile(`^\P{C}+$`)

type User struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if !loginRegex.MatchString(u.Login) {
		return errors.New("invalid Login")
	}
	if !nameRegex.MatchString(u.Name) {
		return errors.New("invalid Name")
	}
	return nil
}
