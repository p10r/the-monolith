package statistics

import "github.com/p10r/pedro/serve/domain"

type Scraper interface {
	GetStatsFor(dm domain.FinishedMatches) (
		matched domain.FinishedMatches,
		notFound domain.FinishedMatches,
		err error,
	)
}
