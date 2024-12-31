package testutil

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"testing"
)

func NewDiscordServer(t *testing.T, logger *slog.Logger) http.HandlerFunc {
	t.Helper()
	log := logger.With(slog.String("adapter", "discord_fake"))

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			log.Info("DiscordServer: Received invalid request")
			w.WriteHeader(400)
			return
		}

		log.Info("DiscordServer: Received request")

		w.WriteHeader(204)
	}
}

func NewFlashscoreServer(t *testing.T, apiKey string) http.HandlerFunc {
	t.Helper()

	//nolint
	//https://flashscore.p.rapidapi.com/v1/events/list?locale=en_GB&timezone=-4&sport_id=12&indent_days=0
	return func(w http.ResponseWriter, r *http.Request) {
		apiKeyHeader := r.Header.Get("X-RapidAPI-Key")
		if apiKeyHeader != apiKey {
			w.WriteHeader(400)
			t.Fatalf("Flashscore Server: X-RapidAPI-Key does not match. "+
				"\n\t\tGot: %v "+
				"\n\t\tWant: %v", apiKeyHeader, apiKey)
			return
		}

		body, err := json.Marshal(FlashscoreRes(t))
		if err != nil {
			t.Fatal("could not marshall JSON")
		}
		// TODO: check for X-RapidAPI-Host and X-RapidAPI-Key
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(body)
		if err != nil {
			t.Fatalf("could not set response: %v", err)
		}
	}
}

func NewPlusLigaServer(t *testing.T, resBody []byte) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resBody)
	}
}

func NewSuperLegaServer(t *testing.T, resBody []byte) http.HandlerFunc {
	t.Helper()

	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(resBody)
	}
}
