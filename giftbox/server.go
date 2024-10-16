package giftbox

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

func NewServer(
	store *sync.Map,
	newUUID func() (string, error),
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /gifts/sweets", panicMiddleware(handleAddSweet(store, newUUID)))
	mux.Handle("POST /gifts/redeem", panicMiddleware(handleRedeemGift(store)))
	return mux
}

type Gift struct {
	ID       string
	Type     string
	Redeemed bool
}

type GiftAddedRes struct {
	ID string `json:"id"`
}

func handleAddSweet(store *sync.Map, newUUID func() (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := newUUID()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		store.Store(id, Gift{ID: id, Type: "SWEET", Redeemed: false})

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		res := GiftAddedRes{ID: id}
		//nolint:errcheck
		json.NewEncoder(w).Encode(res)
	}
}

func handleRedeemGift(store *sync.Map) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("pling")
		giftID := r.URL.Query().Get("id")
		if giftID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		v, ok := store.Load(giftID)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		entry := v.(Gift)
		if entry.Redeemed {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		entry.Redeemed = true
		store.Store(entry.ID, entry)

		w.WriteHeader(http.StatusOK)
	}
}
