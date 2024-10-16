package giftbox

import (
	"fmt"
	"log"
	"net/http"
)

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
