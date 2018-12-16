package core

import "time"

type User struct {
	ID        int64
	Login     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
