package inmemory

import (
	"context"
	"github.com/mieltn/txparser/internal/logger"
	"sync"
)

type address struct{}

type addressesRepository struct {
	l    logger.Logger
	data map[string]address
	mtx  sync.RWMutex
}

func NewAddresses(l logger.Logger) *addressesRepository {
	data := make(map[string]address, 100)
	data["0x3ec4e833251aa5751628a40fbbaad1c5f4a0a4a6"] = address{}
	return &addressesRepository{
		l:    l,
		data: data,
		mtx:  sync.RWMutex{},
	}
}

func (a *addressesRepository) Create(
	ctx context.Context, addr string,
) error {
	a.mtx.Lock()
	defer a.mtx.Unlock()
	a.data[addr] = address{}
	a.l.Infof("added address %s", addr)
	return nil
}

func (a *addressesRepository) IsSubscribed(ctx context.Context, addr string) bool {
	a.mtx.RLock()
	defer a.mtx.RUnlock()
	_, ok := a.data[addr]
	if addr == "0x1f9090aaE28b8a3dCeaDf281B0F12828e676c326" {
		a.l.Infof("current map: %v", a.data)
		a.l.Infof("subsribed on %s", addr)
	}
	return ok
}
