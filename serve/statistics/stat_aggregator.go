package statistics

import (
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/serve/domain"
	"log/slog"
	"net/http"
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
	client *http.Client,
) *Aggregator {
	return &Aggregator{
		newPlusLiga(plusLigaBaseUrl, client),
		newSuperLegaScraper(superLegaBaseUrl, client),
		log,
	}
}

func (a *Aggregator) GetItalianMenStats() domain.StatSheets {
	stats, err := a.superLega.GetStats()
	if err != nil {
		a.log.Error(l.Error("Stats - Italy err: %w", err))
	}
	if stats == nil {
		return domain.StatSheets{}
	}

	return stats.toDomain()
}

func (a *Aggregator) GetPolishMenStats() domain.StatSheets {
	stats, err := a.plusLiga.GetStats()
	if err != nil {
		a.log.Error(l.Error("Stats - Poland err: %w", err))
	}
	if stats == nil {
		return domain.StatSheets{}
	}

	return stats.toDomain()
}
