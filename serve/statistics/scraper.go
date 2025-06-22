package statistics

import "github.com/p10r/monolith/serve/domain"

type Scraper interface {
	GetStatsFor(dm domain.Matches) (
		matched domain.Matches,
		notFound domain.Matches,
		err error,
	)
}
