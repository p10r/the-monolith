package env

import (
	"context"
	"fmt"
	"github.com/p10r/monolith/giftbox"
	"github.com/p10r/monolith/pkg/sqlite"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
)

type InMemory struct {
	Server       http.Handler
	DB           *sqlite.DB
	Repo         *giftbox.GiftRepository
	EventMonitor *InMemoryEventMonitor
	IdGen        func() (string, error)
	apiKey       string
}

type InMemoryEventMonitor struct {
	Events []giftbox.Event
}

func (m *InMemoryEventMonitor) Track(e giftbox.Event) {
	m.Events = append(m.Events, e)
}

func NewInMemoryEnv(t *testing.T, initialID int32, apiKey string) *InMemory {
	ctx := context.Background()
	idGen := func() (string, error) {
		current := atomic.AddInt32(&initialID, 1)
		return fmt.Sprint(current), nil
	}
	db := sqlite.MustOpenDB(t)
	monitor := &InMemoryEventMonitor{Events: make([]giftbox.Event, 0)}

	server, err := giftbox.NewServer(ctx, db, idGen, apiKey, monitor)
	if err != nil {
		t.Fatal(err)
	}
	return &InMemory{
		Server:       server,
		DB:           db,
		Repo:         giftbox.NewGiftRepository(db),
		EventMonitor: monitor,
		IdGen:        idGen,
		apiKey:       apiKey,
	}
}

func (env *InMemory) Events() []giftbox.Event {
	return env.EventMonitor.Events
}

func (env *InMemory) FindInDB(t *testing.T, id giftbox.GiftID) giftbox.Gift {
	gifts, err := env.Repo.All(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for _, gift := range gifts {
		if gift.ID == id {
			return gift
		}
	}
	return giftbox.Gift{}
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
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}

func (env *InMemory) ListAllGifts() *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/gifts", nil)
	req.Header.Set(giftbox.HeaderApiKey, env.apiKey)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}

func (env *InMemory) ListAllPendingGifts() *httptest.ResponseRecorder {
	req := httptest.NewRequest("GET", "/gifts?pending-only=true", nil)
	req.Header.Set(giftbox.HeaderApiKey, env.apiKey)
	w := httptest.NewRecorder()

	env.Server.ServeHTTP(w, req)

	return w
}
