package env

import (
	"fmt"
	"github.com/p10r/pedro/giftbox"
	"github.com/p10r/pedro/pkg/sqlite"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
)

type InMemory struct {
	Server http.Handler
	DB     *sqlite.DB
	Store  *sync.Map
	IdGen  func() (string, error)
}

func NewInMemoryEnv(t *testing.T, initialID int32) *InMemory {
	var store sync.Map
	idGen := func() (string, error) {
		current := atomic.AddInt32(&initialID, 1)
		return fmt.Sprint(current), nil
	}

	return &InMemory{
		Server: giftbox.NewServer(&store, idGen),
		Store:  &store,
		IdGen:  idGen,
	}
}

func (env *InMemory) CheckStoreFor(id string) (giftbox.Gift, bool) {
	value, ok := env.Store.Load(id)
	return value.(giftbox.Gift), ok
}

func (env *InMemory) AddSweet() *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/gifts/sweets", nil)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}

func (env *InMemory) RedeemGift(id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/gifts/redeem?id="+id, nil)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}
