package serve

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alecthomas/assert/v2"
	"github.com/p10r/monolith/pkg/l"
)

func TestFlashscoreResponse(t *testing.T) {
	var response flashscoreResponse
	err := json.Unmarshal(mustReadFile(t, "testdata/flashscore-res.json"), &response)

	assert.NoError(t, err)
	assert.Equal(t, 6, len(response.Leagues))
}

func TestFlashscoreMapping(t *testing.T) {
	response := flashscoreResponse{
		Leagues: []flashscoreLeague{{
			Name: "Germany: 1. Bundesliga",
			Events: []flashscoreEvent{{
				HomeName:         "Berlin",
				AwayName:         "Haching",
				StartTime:        1234,
				HomeScoreCurrent: "3",
				AwayScoreCurrent: "1",
				Stage:            "FINISHED",
			}},
		}},
	}

	expected := Matches{{
		HomeName:         "Berlin",
		AwayName:         "Haching",
		StartTime:        1234,
		FlashscoreName:   "Germany: 1. Bundesliga",
		Country:          "Germany",
		League:           "1. Bundesliga",
		Stage:            "FINISHED",
		HomeScoreCurrent: 3,
		AwayScoreCurrent: 1,
	}}

	assert.Equal(t, expected, response.toMatches())
}

func TestFetchMatches(t *testing.T) {
	apiKey := "1234"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-RapidAPI-Key") != apiKey {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, _ = w.Write(mustReadFile(t, "testdata/flashscore-res.json"))
	}))
	defer server.Close()

	app := New(server.URL, apiKey, "https://discord.invalid", l.NewTextLogger(), "", "")

	matches, err := app.fetchMatches()

	assert.NoError(t, err)
	assert.True(t, len(matches) > 0)
}

func TestFlashscoreFetchError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	app := New(server.URL, "api-key", "https://discord.invalid", l.NewTextLogger(), "", "")

	_, err := app.fetchMatches()
	assert.Error(t, err)
}
