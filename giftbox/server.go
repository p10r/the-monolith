package giftbox

import (
	"context"
	"encoding/json"
	"github.com/p10r/pedro/pkg/sqlite"
	"log"
	"net/http"
)

func NewServer(
	ctx context.Context,
	conn *sqlite.DB,
	newUUID func() (string, error),
	apiKey string,
) http.Handler {
	if apiKey == "" {
		log.Fatal("no api key provided")
	}

	repo := NewGiftRepository(conn)
	idMiddleware := func(next http.Handler) http.Handler {
		return giftIdMiddleware(ctx, newUUID, next)
	}
	auth := func(next http.Handler) http.Handler {
		return authMiddleWare(apiKey, next)
	}

	mux := http.NewServeMux()
	// TODO add uuid in separate middleware
	mux.Handle("POST /gifts/sweets", auth(idMiddleware(handleAddSweet(repo))))
	mux.Handle("POST /gifts/wishes", auth(idMiddleware(handleAddWish(repo))))
	mux.Handle("POST /gifts/images", auth(idMiddleware(handleAddImage(repo))))
	// Using a GET here as it's called via QR code
	mux.Handle("GET /gifts/redeem", handleRedeemGift(repo))

	return panicMiddleware(mux)
}

type GiftAddedRes struct {
	ID string `json:"id"`
}

func handleAddSweet(
	repo *GiftRepository,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO check ID
		id := r.Context().Value(ctxGiftID).(GiftID)
		gift := Gift{ID: id, Type: TypeSweet, Redeemed: false}

		err := repo.Save(context.Background(), gift)
		if err != nil {
			log.Printf("err when writing to db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		res := GiftAddedRes{ID: id.String()}
		//nolint:errcheck
		json.NewEncoder(w).Encode(res)
	}
}

func handleAddWish(
	repo *GiftRepository,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ctxGiftID).(GiftID)
		gift := Gift{ID: id, Type: TypeWish, Redeemed: false}

		err := repo.Save(context.Background(), gift)
		if err != nil {
			log.Printf("err when writing to db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		res := GiftAddedRes{ID: id.String()}
		//nolint:errcheck
		json.NewEncoder(w).Encode(res)
	}
}

func handleAddImage(
	repo *GiftRepository,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		imgUrl := r.URL.Query().Get("url")
		if imgUrl == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		id := r.Context().Value(ctxGiftID).(GiftID)
		gift := Gift{ID: id, Type: TypeImage, Redeemed: false, ImageUrl: imgUrl}

		err := repo.Save(context.Background(), gift)
		if err != nil {
			log.Printf("err when writing to db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		res := GiftAddedRes{ID: id.String()}
		//nolint:errcheck
		json.NewEncoder(w).Encode(res)
	}
}

func handleRedeemGift(repo *GiftRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqId := r.URL.Query().Get("id")
		if reqId == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		gifts, err := repo.All(context.Background())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
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

		_, err = repo.SetRedeemedFlag(context.Background(), gift.ID.String(), true)
		if err != nil {
			log.Printf("err when writing to db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if gift.Type == TypeImage {
			w.Header().Set("Location", gift.ImageUrl)
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
