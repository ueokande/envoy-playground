package blob

import uuid "github.com/satori/go.uuid"

func NewID() string {
	return uuid.NewV4().String()
}
