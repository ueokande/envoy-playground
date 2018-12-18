package mock

import (
	"sync"

	core "github.com/ueokande/envoy-playground"
	"github.com/ueokande/envoy-playground/db"
)

func New() db.Interface {
	return &impl{
		users:      make(map[string]core.User),
		userPhotos: make(map[string]string),
	}
}

type impl struct {
	m sync.Mutex

	users      map[string]core.User
	userPhotos map[string]string
}

func (i *impl) Close() error {
	return nil
}
