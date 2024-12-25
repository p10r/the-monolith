package statistics

import (
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/serve/domain"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"
)

type Aggregator struct {
	plusLiga *plusLigaScraper
	log      *slog.Logger
}

func NewAggregator(plusLigaBaseUrl string, log *slog.Logger) *Aggregator {
	c := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 10 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return &Aggregator{newPlusLiga(plusLigaBaseUrl, c), log}
}

func (a *Aggregator) EnrichMatches(matches domain.FinishedMatches) domain.FinishedMatches {
	matchesMap := map[matchKey]domain.FinishedMatch{}
	for _, match := range matches {
		matchesMap[newMatchKey(match.HomeName, match.AwayName)] = match
	}

	// Get all PlusLiga Matches
	var plMatches domain.FinishedMatches
	for _, match := range matches {
		if strings.ToLower(match.Country) == "poland" {
			plMatches = append(plMatches, match)
		}
	}

	// Get stats for PlusLiga matches
	found, notFound, err := a.plusLiga.GetStatsFor(plMatches)
	if err != nil {
		a.log.Error(l.Error("Plusliga err: %w", err))
	}
	for _, match := range notFound {
		a.log.Error("Not found on PlusLiga website: %s-%s", match.HomeName, match.AwayName)
	}

	// Overwrite statsUrl of domain.Match
	for _, foundMatch := range found {
		matchesMap[newMatchKey(foundMatch.HomeName, foundMatch.AwayName)] = foundMatch
	}

	// map back to slice
	var unwrapped domain.FinishedMatches
	for _, match := range matchesMap {
		unwrapped = append(unwrapped, match)
	}
	return unwrapped
}
