package scraper

import (
	"github.com/alecthomas/assert/v2"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/testutil"
	"strings"
	"testing"
)

func TestPlusLigaScraper(t *testing.T) {
	plusLigaBaseUrl := "plusliga-url"
	plusLiga := PlusLiga{baseUrl: plusLigaBaseUrl}

	t.Run("scrapes matches", func(t *testing.T) {
		f := testutil.MustReadFile(t, "../testdata/scraper/plusliga-game-day.html")
		res := testutil.SomeRes(f)

		stats, err := plusLiga.ParseStats(res)
		expected := []plusLigaMatch{
			{
				homeTeam: "BOGDANKA LUK Lublin",
				awayTeam: "Ślepsk Malow Suwałki",
				statsUrl: "plusliga-url/games/action/show/id/1103632.html",
			},
			{
				homeTeam: "Aluron CMC Warta Zawiercie",
				awayTeam: "Steam Hemarpol Norwid Częstochowa",
				statsUrl: "plusliga-url/games/action/show/id/1103637.html",
			},
			{
				homeTeam: "PSG Stal Nysa",
				awayTeam: "Cuprum Stilon Gorzów",
				statsUrl: "plusliga-url/games/action/show/id/1103635.html",
			},
			{
				homeTeam: "Nowak-Mosty MKS Będzin",
				awayTeam: "Trefl Gdańsk",
				statsUrl: "plusliga-url/games/action/show/id/1103634.html",
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, expected, stats)
	})

	t.Run("scrapes matches from full page", func(t *testing.T) {
		f := testutil.MustReadFile(t, "../testdata/scraper/plusliga.html")
		res := testutil.SomeRes(f)

		stats, err := plusLiga.ParseStats(res)
		assert.NoError(t, err)
		assert.Equal(t, 240, len(stats))

		for _, m := range stats {
			assert.False(t, m.homeTeam == "")
			assert.False(t, m.awayTeam == "")
			assert.False(t, m.statsUrl == "")
			assert.True(t, strings.HasPrefix(m.statsUrl, plusLigaBaseUrl))
		}
	})

	t.Run("maps plusliga matches to Serve Matches", func(t *testing.T) {
		plm := plusLigaMatches{
			{
				homeTeam: "BOGDANKA LUK Lublin",
				awayTeam: "Ślepsk Malow Suwałki",
				statsUrl: "plusliga-url/games/action/show/id/1103632.html",
			},
			//{
			//	homeTeam: "Aluron CMC Warta Zawiercie",
			//	awayTeam: "Steam Hemarpol Norwid Częstochowa",
			//	statsUrl: "plusliga-url/games/action/show/id/1103637.html",
			//},
		}

		domainMatches := domain.Matches{
			domain.Match{
				ID:        123,
				HomeName:  "Lublin",
				AwayName:  "Slepsk Suwalki",
				StartTime: 1,
				Country:   "Poland",
				League:    "PlusLiga",
				StatsUrl:  "",
			},
		}

		expected := domain.Matches{
			domain.Match{
				ID:        123,
				HomeName:  "Lublin",
				AwayName:  "Slepsk Suwalki",
				StartTime: 1,
				Country:   "Poland",
				League:    "PlusLiga",
				StatsUrl:  "plusliga-url/games/action/show/id/1103632.html",
			},
		}

		zipped, _ := plm.ZipWith(domainMatches)

		assert.Equal(t, expected, zipped)
	})
}
