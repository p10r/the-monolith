package giftbox

import (
	"encoding/json"
	"net/http"
	"sync"
)

func NewServer(
	store *sync.Map,
	newUUID func() (string, error),
) http.Handler {
	mux := http.NewServeMux()
	var handler http.Handler = mux

	mux.Handle("POST /gifts/sweets", handleAddSweet(store, newUUID))

	return handler
}

type GiftAddedRes struct {
	ID string `json:"id"`
}

func handleAddSweet(store *sync.Map, newUUID func() (string, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := newUUID()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		store.Store(id, "sweet")

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		res := GiftAddedRes{ID: id}
		//nolint:errcheck
		json.NewEncoder(w).Encode(res)
	}
}
