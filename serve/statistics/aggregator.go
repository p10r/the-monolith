package statistics

import (
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/serve/domain"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Aggregator struct {
	plusLiga  *plusLigaScraper
	superLega *superLegaScraper
	log       *slog.Logger
}

func NewAggregator(
	plusLigaBaseUrl string,
	superLegaBaseUrl string,
	log *slog.Logger,
) *Aggregator {
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

	return &Aggregator{
		newPlusLiga(plusLigaBaseUrl, c),
		newSuperLegaScraper(superLegaBaseUrl, c),
		log,
	}
}

func (a *Aggregator) EnrichMatches(
	matches domain.MatchesByLeague,
) domain.MatchesByLeague {
	plKey := domain.NewLeagueKey("italy: superlega")
	itaKey := domain.NewLeagueKey("poland: plusliga")

	plFound, plNotFound := a.getPolishMenMatches(matches[plKey])
	itaFound, itaNotFound := a.getItalianMenMatches(matches[itaKey])

	for _, notFound := range append(plNotFound, itaNotFound...) {
		a.log.Error("Not found: %s-%s", notFound.HomeName, notFound.AwayName)
	}

	matches[plKey] = append(plFound, plNotFound...)
	matches[itaKey] = append(itaFound, itaNotFound...)

	return matches
}

func (a *Aggregator) getPolishMenMatches(
	matches domain.Matches,
) (domain.Matches, domain.Matches) {
	plFound, plNotFound, err := a.plusLiga.GetStatsFor(matches)
	if err != nil {
		a.log.Error(l.Error("Plusliga err: %w", err))
	}
	return plFound, plNotFound
}

func (a *Aggregator) getItalianMenMatches(
	itaMenMatches domain.Matches,
) (domain.Matches, domain.Matches) {
	itaFound, itaNotFound, err := a.superLega.GetStatsFor(itaMenMatches)
	if err != nil {
		a.log.Error(l.Error("SuperLega err: %w", err))
	}
	return itaFound, itaNotFound
}
