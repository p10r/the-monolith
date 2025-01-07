package statistics

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/testutil"
	"net/http"
	"testing"
)

func NewFixture(tripper testutil.RoundTripFunc) domain.Statistics {
	return NewAggregator(
		"/plusliga",
		"/superlega",
		l.NewTextLogger(),
		testutil.NewTestClient(tripper),
	)
}

func TestStats(t *testing.T) {
	t.Run("returns stats for Italy's SuperLega", func(t *testing.T) {
		statistics := NewFixture(func(req *http.Request) *http.Response {
			if req.URL.String() == "/superlega/calendario/?lang=en" {
				//nolint:lll
				return testutil.OkRes(testutil.MustReadFile(t, "../testdata/statistics/superlega-italy-m.html"))
			}
			panic(fmt.Sprintf("err, req URL was: %s", req.URL.String()))
		})

		//nolint:lll
		expected := domain.StatSheets{
			{League: "italy: superlega", Home: "Verona", Away: "Perugia", Url: "http://www.legavolley.it/match/38270"},
			{League: "italy: superlega", Home: "Modena", Away: "Piacenza", Url: "http://www.legavolley.it/match/38271"},
			{League: "italy: superlega", Home: "Milano", Away: "Taranto", Url: "http://www.legavolley.it/match/38274"},
			{League: "italy: superlega", Home: "Monza", Away: "Grottazzolina", Url: "http://www.legavolley.it/match/38275"},
			{League: "italy: superlega", Home: "Taranto", Away: "Trentino", Url: "http://www.legavolley.it/match/38276"},
			{League: "italy: superlega", Home: "Piacenza", Away: "Monza", Url: "http://www.legavolley.it/match/38277"},
			{League: "italy: superlega", Home: "Lube Civitanova", Away: "Milano", Url: "http://www.legavolley.it/match/38278"},
			{League: "italy: superlega", Home: "Cisterna", Away: "Verona", Url: "http://www.legavolley.it/match/38279"},
			{League: "italy: superlega", Home: "Grottazzolina", Away: "Modena", Url: "http://www.legavolley.it/match/38280"},
			{League: "italy: superlega", Home: "Perugia", Away: "Padova", Url: "http://www.legavolley.it/match/38281"},
			{League: "italy: superlega", Home: "Monza", Away: "Lube Civitanova", Url: "http://www.legavolley.it/match/38282"},
			{League: "italy: superlega", Home: "Trentino", Away: "Milano", Url: "http://www.legavolley.it/match/38283"},
			{League: "italy: superlega", Home: "Perugia", Away: "Modena", Url: "http://www.legavolley.it/match/38284"},
			{League: "italy: superlega", Home: "Piacenza", Away: "Cisterna", Url: "http://www.legavolley.it/match/38285"},
			{League: "italy: superlega", Home: "Padova", Away: "Taranto", Url: "http://www.legavolley.it/match/38286"},
			{League: "italy: superlega", Home: "Verona", Away: "Grottazzolina", Url: "http://www.legavolley.it/match/38287"},
		}

		assert.Equal(t, expected, statistics.GetItalianMenStats())
	})

	t.Run("returns stats for Poland's PlusLiga", func(t *testing.T) {
		statistics := NewFixture(func(req *http.Request) *http.Response {
			if req.URL.String() == "/plusliga/games.html" {
				return testutil.OkRes(testutil.MustReadFile(t, "../testdata/statistics/plusliga.html"))
			}
			panic(fmt.Sprintf("err, req URL was: %s", req.URL.String()))
		})

		firstMatch := domain.StatSheet{
			League: "poland: plusliga",
			Home:   "Zawierce",
			Away:   "Lublin",
			Url:    "/plusliga/games/action/show/id/1103460.html",
		}
		secondMatch := domain.StatSheet{
			League: "poland: plusliga",
			Home:   "Projekt Warszawa",
			Away:   "Kedzierzyn-Kozle",
			Url:    "/plusliga/games/action/show/id/1103464.html",
		}

		actual := statistics.GetPolishMenStats()
		assert.Equal(t, 240, len(actual))
		assert.Equal(t, firstMatch, actual[0])
		assert.Equal(t, secondMatch, actual[1])
	})

}
