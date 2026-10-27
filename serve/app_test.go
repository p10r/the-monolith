package serve

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
	"github.com/p10r/monolith/pkg/l"
)

// TestServeEndToEnd runs the whole flow against fake Flashscore and Discord
// servers: fetch matches, keep favourites, split scheduled/finished, format and
// post the Discord messages.
func TestServeEndToEnd(t *testing.T) {
	const apiKey = "test-api-key"

	var (
		mu       sync.Mutex
		received []discordMessage
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /flashscore/v1/events/list", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-RapidAPI-Key") != apiKey {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, _ = w.Write(mustReadFile(t, "testdata/flashscore-res.json"))
	})
	mux.HandleFunc("POST /discord", func(w http.ResponseWriter, r *http.Request) {
		var msg discordMessage
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		mu.Lock()
		received = append(received, msg)
		mu.Unlock()

		w.WriteHeader(http.StatusNoContent)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	app := New(server.URL+"/flashscore", apiKey, server.URL+"/discord", l.NewTextLogger(), "", "")
	app.now = func() time.Time {
		return time.Date(2024, 5, 28, 0, 0, 0, 0, time.UTC)
	}

	scheduled, err := app.importScheduledMatches()
	assert.NoError(t, err)
	assert.Equal(t, 3, len(scheduled))

	err = app.importFinishedMatches()
	assert.NoError(t, err)

	msgByType := messagesByType(t, received)

	upcoming := msgByType["games"]
	assert.Equal(t, "Games for Tuesday, 28 May 2024", upcoming.Content)

	upcomingFields := fieldsByName(upcoming)
	assert.Equal(t, 2, len(upcomingFields))
	assert.True(t,
		strings.Contains(upcomingFields["🇪🇺EUROPE: CHAMPIONS LEAGUE - PLAY OFFS"], "**Trentino - Jastrzebski**"),
		"expected upcoming Europe match in %q", upcomingFields["🇪🇺EUROPE: CHAMPIONS LEAGUE - PLAY OFFS"],
	)
	assert.True(t,
		strings.Contains(upcomingFields["🇺🇸USA: PVF WOMEN"], "**Grand Rapids Rise W - San Diego Mojo W**"),
		"expected upcoming USA match in %q", upcomingFields["🇺🇸USA: PVF WOMEN"],
	)

	finished := msgByType["results"]
	assert.Equal(t, "Results for Tuesday, 28 May 2024", finished.Content)

	finishedFields := fieldsByName(finished)
	assert.Equal(t, 2, len(finishedFields))

	italy := finishedFields["🇮🇹italy: superlega"]
	assert.Equal(t, 6, len(strings.Split(italy, "\n")))
	assert.True(t,
		strings.Contains(italy, "**Monza - Lube Civitanova**\t\t\tScore:\t||3\t:\t2||"),
		"expected finished Italy match in %q", italy,
	)

	europe := finishedFields["🇪🇺europe: champions league women - play offs"]
	assert.True(t,
		strings.Contains(europe, "**Conegliano W - Vero Volley W**\t\t\tScore:\t||3\t:\t2||"),
		"expected finished Europe match in %q", europe,
	)
}

func messagesByType(t *testing.T, messages []discordMessage) map[string]discordMessage {
	t.Helper()
	byType := make(map[string]discordMessage, len(messages))
	for _, msg := range messages {
		switch {
		case strings.HasPrefix(msg.Content, "Games for"):
			byType["games"] = msg
		case strings.HasPrefix(msg.Content, "Results for"):
			byType["results"] = msg
		default:
			t.Fatalf("unexpected discord message: %q", msg.Content)
		}
	}
	assert.Equal(t, 2, len(messages))
	return byType
}

func fieldsByName(msg discordMessage) map[string]string {
	fields := make(map[string]string, len(msg.Embeds[0].Fields))
	for _, field := range msg.Embeds[0].Fields {
		fields[field.Name] = field.Value
	}
	return fields
}
