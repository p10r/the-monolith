package giftbox

import (
	"context"
	"encoding/json"
	"github.com/p10r/pedro/pkg/sqlite"
	"log"
	"net/http"
)

func NewServer(
	conn *sqlite.DB,
	newUUID func() (string, error),
) http.Handler {
	repo := NewGiftRepository(conn)

	mux := http.NewServeMux()
	mux.Handle("POST /gifts/sweets", handleAddSweet(repo, newUUID))
	mux.Handle("POST /gifts/redeem", handleRedeemGift(repo))
	return panicMiddleware(mux)
}

type Gifts []Gift

func (g Gifts) findByID(reqId string) (Gift, bool) {
	giftsByID := make(map[string]Gift)
	for _, gift := range g {
		giftsByID[gift.ID] = gift
	}
	gift, ok := giftsByID[reqId]
	return gift, ok
}

type GiftType string

const (
	TypeSweet GiftType = "SWEET"
	TypeWish  GiftType = "WISH"
)

type Gift struct {
	ID       string
	Type     GiftType
	Redeemed bool
}

type GiftAddedRes struct {
	ID string `json:"id"`
}

func handleAddSweet(
	repo *GiftRepository,
	newUUID func() (string, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := newUUID()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		gift := Gift{ID: id, Type: "SWEET", Redeemed: false}

		err = repo.Save(context.Background(), gift)
		if err != nil {
			log.Printf("err when writing to db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		res := GiftAddedRes{ID: id}
		//nolint:errcheck
		json.NewEncoder(w).Encode(res)
	}
}

func handleRedeemGift(repo *GiftRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("pling")
		reqId := r.URL.Query().Get("id")
		if reqId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		gifts, err := repo.All(context.Background())
		if err != nil {
			return
		}

		gift, ok := gifts.findByID(reqId)
		if !ok {
			log.Printf("gift %s could not be found in db", reqId)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if gift.Redeemed {
			log.Printf("gift %s is already redeemed", gift.ID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = repo.SetRedeemedFlag(context.Background(), gift.ID, true)
		if err != nil {
			log.Printf("err when writing to db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
