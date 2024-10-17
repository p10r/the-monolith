package env

import (
	"context"
	"fmt"
	"github.com/p10r/pedro/giftbox"
	"github.com/p10r/pedro/pkg/sqlite"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

type InMemory struct {
	Server http.Handler
	DB     *sqlite.DB
	Repo   *giftbox.GiftRepository
	IdGen  func() (string, error)
}

func NewInMemoryEnv(t *testing.T, initialID int32) *InMemory {
	idGen := func() (string, error) {
		current := atomic.AddInt32(&initialID, 1)
		return fmt.Sprint(current), nil
	}
	db := sqlite.MustOpenDB(t)

	return &InMemory{
		Server: giftbox.NewServer(db, idGen),
		DB:     db,
		Repo:   giftbox.NewGiftRepository(db),
		IdGen:  idGen,
	}
}

func (env *InMemory) FindInDB(t *testing.T, id string) (giftbox.Gift, bool) {
	gifts, err := env.Repo.All(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, gift := range gifts {
		if gift.ID == id {
			return gift, true
		}
	}
	return giftbox.Gift{}, false
	//value, ok := env.Store.Load(id)
	//return value.(giftbox.Gift), ok
}

func (env *InMemory) AddSweet() *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/gifts/sweets", nil)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}

func (env *InMemory) AddWish() *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/gifts/wishes", nil)
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
