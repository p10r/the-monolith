package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/p10r/monolith/pkg/l"
)

type flashscoreResponse struct {
	Leagues []flashscoreLeague `json:"DATA"`
}

type flashscoreLeague struct {
	Name   string            `json:"NAME"`
	Events []flashscoreEvent `json:"EVENTS"`
}

type flashscoreEvent struct {
	HomeName         string `json:"HOME_NAME"`
	AwayName         string `json:"AWAY_NAME"`
	StartTime        int64  `json:"START_TIME"`
	HomeScoreCurrent string `json:"HOME_SCORE_CURRENT"`
	AwayScoreCurrent string `json:"AWAY_SCORE_CURRENT"`
	Stage            string `json:"STAGE"`
}

func (s *Serve) fetchMatches() (Matches, error) {
	url := s.flashscoreBaseURL + "/v1/events/list?locale=en_GB&timezone=-4&sport_id=12&indent_days=0"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create flashscore request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-RapidAPI-Host", "flashscore.p.rapidapi.com")
	req.Header.Set("X-RapidAPI-Key", s.flashscoreAPIKey)

	res, err := s.doWithRetries(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusForbidden {
		return nil, errors.New("flashscore: forbidden - wrong API key?")
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flashscore request failed: %s", res.Status)
	}

	var response flashscoreResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode flashscore response: %w", err)
	}

	return response.toMatches(), nil
}

func (s *Serve) doWithRetries(req *http.Request) (*http.Response, error) {
	const maxRetries = 5

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		res, err := s.httpClient.Do(req)
		if err == nil {
			return res, nil
		}
		lastErr = err
		s.log.Error(l.Error("flashscore request failed", err))
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("flashscore request failed after %d attempts: %w", maxRetries, lastErr)
}

func (r flashscoreResponse) toMatches() Matches {
	var matches Matches
	for _, league := range r.Leagues {
		country, leagueName, _ := strings.Cut(league.Name, ": ")
		for _, event := range league.Events {
			home, _ := strconv.Atoi(event.HomeScoreCurrent)
			away, _ := strconv.Atoi(event.AwayScoreCurrent)

			matches = append(matches, Match{
				HomeName:         event.HomeName,
				AwayName:         event.AwayName,
				StartTime:        event.StartTime,
				FlashscoreName:   league.Name,
				Country:          country,
				League:           leagueName,
				Stage:            event.Stage,
				HomeScoreCurrent: home,
				AwayScoreCurrent: away,
			})
		}
	}
	return matches
}
