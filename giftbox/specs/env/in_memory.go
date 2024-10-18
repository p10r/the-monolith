package env

import (
	"context"
	"fmt"
	"github.com/p10r/pedro/giftbox"
	"github.com/p10r/pedro/pkg/sqlite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
)

type InMemory struct {
	Server http.Handler
	DB     *sqlite.DB
	Repo   *giftbox.GiftRepository
	IdGen  func() (string, error)
	apiKey string
}

func NewInMemoryEnv(t *testing.T, initialID int32, apiKey string) *InMemory {
	ctx := context.Background()
	idGen := func() (string, error) {
		current := atomic.AddInt32(&initialID, 1)
		return fmt.Sprint(current), nil
	}
	db := sqlite.MustOpenDB(t)

	return &InMemory{
		Server: giftbox.NewServer(ctx, db, idGen, apiKey),
		DB:     db,
		Repo:   giftbox.NewGiftRepository(db),
		IdGen:  idGen,
		apiKey: apiKey,
	}
}

func (env *InMemory) FindInDB(t *testing.T, id giftbox.GiftID) (giftbox.Gift, bool) {
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
	req.Header.Set(giftbox.HeaderApiKey, env.apiKey)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}

func (env *InMemory) AddWish() *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/gifts/wishes", nil)
	req.Header.Set(giftbox.HeaderApiKey, env.apiKey)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}

func (env *InMemory) AddImage(imageUrl string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("POST", "/gifts/images?url="+url.QueryEscape(imageUrl), nil)
	req.Header.Set(giftbox.HeaderApiKey, env.apiKey)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}

func (env *InMemory) RedeemGift(id string) *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/gifts/redeem?id="+id, nil)
	req.Header.Set(giftbox.HeaderApiKey, env.apiKey)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}
