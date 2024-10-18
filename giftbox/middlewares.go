package giftbox

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const ctxGiftID GiftID = "giftId"

// giftIdMiddleware adds a generated UUID to the request's context
func giftIdMiddleware(
	ctx context.Context,
	newUUID func() (string, error),
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := newUUID()
		if err != nil {
			log.Printf("Failed to generate new UUID: %v", err)
			http.Error(w, "could not generate ID", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, ctxGiftID, GiftID(id))))
	})
}

func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				log.Println(fmt.Errorf("panic middleware: caught err %v", err)) // TODO slog
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			}

		}()
		next.ServeHTTP(w, r)
	})
}
