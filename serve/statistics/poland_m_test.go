package statistics

import (
	"github.com/alecthomas/assert/v2"
	"github.com/p10r/monolith/serve/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPlusLigaScraper(t *testing.T) {
	plusLigaBaseUrl := "https://www.plusliga.pl"
	plusLiga := plusLigaScraper{baseUrl: plusLigaBaseUrl}

	t.Run("scrapes matches", func(t *testing.T) {
		f := testutil.MustReadFile(t, "../testdata/statistics/plusliga-game-day.html")
		res := testutil.OkRes(f)

		stats, err := plusLiga.parseStats(res)
		expected := plusLigaMatches{
			{
				homeTeam: "BOGDANKA LUK Lublin",
				awayTeam: "Ślepsk Malow Suwałki",
				statsUrl: "https://www.plusliga.pl/games/action/show/id/1103632.html",
			},
			{
				homeTeam: "Aluron CMC Warta Zawiercie",
				awayTeam: "Steam Hemarpol Norwid Częstochowa",
				statsUrl: "https://www.plusliga.pl/games/action/show/id/1103637.html",
			},
			{
				homeTeam: "PSG Stal Nysa",
				awayTeam: "Cuprum Stilon Gorzów",
				statsUrl: "https://www.plusliga.pl/games/action/show/id/1103635.html",
			},
			{
				homeTeam: "Nowak-Mosty MKS Będzin",
				awayTeam: "Trefl Gdańsk",
				statsUrl: "https://www.plusliga.pl/games/action/show/id/1103634.html",
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, stats)
	})

	t.Run("scrapes matches from full page", func(t *testing.T) {
		f := testutil.MustReadFile(t, "../testdata/statistics/plusliga.html")
		res := testutil.OkRes(f)

		stats, err := plusLiga.parseStats(res)
		assert.NoError(t, err)
		assert.Equal(t, 240, len(stats))

		for _, m := range stats {
			assert.False(t, m.homeTeam == "")
			assert.False(t, m.awayTeam == "")
			assert.False(t, m.statsUrl == "")
			assert.True(t, strings.HasPrefix(m.statsUrl, plusLigaBaseUrl))
		}
	})

	t.Run("round-trip", func(t *testing.T) {
		f := testutil.MustReadFile(t, "../testdata/statistics/plusliga.html")
		plusLigaSite := httptest.NewServer(testutil.NewPlusLigaServer(t, f))
		defer plusLigaSite.Close()

		plusLiga := newPlusLiga(plusLigaSite.URL, &http.Client{})

		urls, err := plusLiga.GetStats()
		assert.NoError(t, err)
		assert.Equal(t, 240, len(urls))
	})

	t.Run("run against prod", func(t *testing.T) {
		t.Skip()

		plusLiga := newPlusLiga("https://www.plusliga.pl", &http.Client{})

		urls, err := plusLiga.GetStats()
		assert.NoError(t, err)
		assert.Equal(t, 240, len(urls))
	})

}
