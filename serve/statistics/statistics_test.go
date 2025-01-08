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

		expected := domain.StatSheets{
			{Home: "Verona", Away: "Perugia", Url: "http://www.legavolley.it/match/38270"},
			{Home: "Modena", Away: "Piacenza", Url: "http://www.legavolley.it/match/38271"},
			{Home: "Milano", Away: "Taranto", Url: "http://www.legavolley.it/match/38274"},
			{Home: "Monza", Away: "Grottazzolina", Url: "http://www.legavolley.it/match/38275"},
			{Home: "Taranto", Away: "Trentino", Url: "http://www.legavolley.it/match/38276"},
			{Home: "Piacenza", Away: "Monza", Url: "http://www.legavolley.it/match/38277"},
			{Home: "Lube Civitanova", Away: "Milano", Url: "http://www.legavolley.it/match/38278"},
			{Home: "Cisterna", Away: "Verona", Url: "http://www.legavolley.it/match/38279"},
			{Home: "Grottazzolina", Away: "Modena", Url: "http://www.legavolley.it/match/38280"},
			{Home: "Perugia", Away: "Padova", Url: "http://www.legavolley.it/match/38281"},
			{Home: "Monza", Away: "Lube Civitanova", Url: "http://www.legavolley.it/match/38282"},
			{Home: "Trentino", Away: "Milano", Url: "http://www.legavolley.it/match/38283"},
			{Home: "Perugia", Away: "Modena", Url: "http://www.legavolley.it/match/38284"},
			{Home: "Piacenza", Away: "Cisterna", Url: "http://www.legavolley.it/match/38285"},
			{Home: "Padova", Away: "Taranto", Url: "http://www.legavolley.it/match/38286"},
			{Home: "Verona", Away: "Grottazzolina", Url: "http://www.legavolley.it/match/38287"},
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
			Home: "Zawierce",
			Away: "Lublin",
			Url:  "/plusliga/games/action/show/id/1103460.html",
		}
		secondMatch := domain.StatSheet{
			Home: "Projekt Warszawa",
			Away: "Kedzierzyn-Kozle",
			Url:  "/plusliga/games/action/show/id/1103464.html",
		}

		actual := statistics.GetPolishMenStats()
		assert.Equal(t, 240, len(actual))
		assert.Equal(t, firstMatch, actual[0])
		assert.Equal(t, secondMatch, actual[1])
	})

}
