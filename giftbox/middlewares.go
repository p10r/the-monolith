package giftbox

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

const ctxGiftID GiftID = "giftId"
const HeaderApiKey = "X-Gift-Box-Api-Key" // #nosec G101

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

func authMiddleWare(apiKey string, monitor EventMonitor, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqApiKey := r.Header.Get(HeaderApiKey)

		if reqApiKey != apiKey {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			monitor.Track(IllegalAccessEvent{r.URL.String(), string(body)})

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
