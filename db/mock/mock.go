package mock

import (
	"sync"

	core "github.com/ueokande/envoy-playground"
	"github.com/ueokande/envoy-playground/db"
)

func New() db.Interface {
	return &impl{
		users: make(map[string]core.User),
	}
}

type impl struct {
	m sync.Mutex

	users map[string]core.User
}

func (i *impl) Close() error {
	return nil
}
