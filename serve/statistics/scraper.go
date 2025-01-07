package statistics

import "github.com/p10r/pedro/serve/domain"

type Scraper interface {
	GetStatsFor(dm domain.Matches) (
		matched domain.Matches,
		notFound domain.Matches,
		err error,
	)
}
