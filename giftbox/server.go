package giftbox

import (
	"context"
	"encoding/json"
	"github.com/p10r/pedro/pkg/sqlite"
	"html/template"
	"log"
	"net/http"
)

const (
	sweetsGif = "https://media1.tenor.com/m/M3p9DCrC7OkAAAAC/christmas-dinner-sweets.gif"
	//nolint:lll
	wishGif = "https://media1.tenor.com/m/CRN0ZkGmuLkAAAAC/your-wish-is-my-command-jeremy-reynolds.gif"
)

func NewServer(
	ctx context.Context,
	conn *sqlite.DB,
	newUUID func() (string, error),
	apiKey string,
	monitor EventMonitor,
	templateDir string,
) (http.Handler, error) {
	if apiKey == "" {
		log.Fatal("no api key provided")
	}

	tmpl, err := template.ParseFiles(templateDir + "gift-redeemed.html")
	if err != nil {
		return nil, err
	}

	repo := NewGiftRepository(conn)
	idMiddleware := func(next http.Handler) http.Handler {
		return giftIdMiddleware(ctx, newUUID, next)
	}
	auth := func(next http.Handler) http.Handler {
		return authMiddleWare(apiKey, monitor, next)
	}

	mux := http.NewServeMux()
	mux.Handle("POST /gifts/sweets", auth(idMiddleware(handleAddSweet(repo))))
	mux.Handle("POST /gifts/wishes", auth(idMiddleware(handleAddWish(repo))))
	mux.Handle("POST /gifts/images", auth(idMiddleware(handleAddImage(repo))))
	mux.Handle("GET /gifts", auth(idMiddleware(handleListAllGifts(repo))))
	// Using a GET here as it's called via QR code
	mux.Handle("GET /gifts/redeem", handleRedeemGift(repo, monitor, tmpl))

	return panicMiddleware(mux), nil
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

func handleRedeemGift(
	repo *GiftRepository,
	monitor EventMonitor,
	tmpl *template.Template,
) http.HandlerFunc {
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
			monitor.Track(NotFoundEvent{ID: reqId})
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if gift.Redeemed {
			monitor.Track(AlreadyRedeemedEvent{gift.ID, gift.Type})
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = repo.SetRedeemedFlag(context.Background(), gift.ID.String(), true)
		if err != nil {
			log.Printf("err when writing to db: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		monitor.Track(RedeemedEvent{
			ID:   gift.ID,
			Type: gift.Type,
		})

		if gift.Type == TypeImage {
			w.Header().Set("Location", gift.ImageUrl)
			w.WriteHeader(http.StatusSeeOther)
			return
		}

		if gift.Type == TypeSweet {
			_ = tmpl.Execute(w, sweetsGif)
			return
		}

		_ = tmpl.Execute(w, wishGif)
	}
}

type AllGiftsRes struct {
	Gifts Gifts `json:"gifts"`
}

func handleListAllGifts(repo *GiftRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pendingOnly bool
		if r.URL.Query().Get("pending-only") == "true" {
			pendingOnly = true
		} else {
			pendingOnly = false
		}

		gifts, err := repo.All(context.Background())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !pendingOnly {
			w.WriteHeader(http.StatusOK)
			//nolint:errcheck
			json.NewEncoder(w).Encode(AllGiftsRes{Gifts: gifts})
			return
		}

		var outstandingGifts Gifts
		for _, gift := range gifts {
			if !gift.Redeemed {
				outstandingGifts = append(outstandingGifts, gift)
			}
		}

		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		json.NewEncoder(w).Encode(AllGiftsRes{Gifts: outstandingGifts})
	}
}
