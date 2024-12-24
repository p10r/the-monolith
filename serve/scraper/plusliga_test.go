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

		stats, err := plusLiga.parseStats(res)
		expected := plusLigaMatches{
			newMatchKey("BOGDANKA LUK Lublin", "Ślepsk Malow Suwałki"): {
				homeTeam: "BOGDANKA LUK Lublin",
				awayTeam: "Ślepsk Malow Suwałki",
				statsUrl: "plusliga-url/games/action/show/id/1103632.html",
			},
			newMatchKey("Aluron CMC Warta Zawiercie", "Steam Hemarpol Norwid Częstochowa"): {
				homeTeam: "Aluron CMC Warta Zawiercie",
				awayTeam: "Steam Hemarpol Norwid Częstochowa",
				statsUrl: "plusliga-url/games/action/show/id/1103637.html",
			},
			newMatchKey("PSG Stal Nysa", "Cuprum Stilon Gorzów"): {
				homeTeam: "PSG Stal Nysa",
				awayTeam: "Cuprum Stilon Gorzów",
				statsUrl: "plusliga-url/games/action/show/id/1103635.html",
			},
			newMatchKey("Nowak-Mosty MKS Będzin", "Trefl Gdańsk"): {
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

	t.Run("maps PlusLiga matches to Serve Matches", func(t *testing.T) {
		match1 := plusLigaMatch{
			homeTeam: "BOGDANKA LUK Lublin",
			awayTeam: "Ślepsk Malow Suwałki",
			statsUrl: "plusliga-url/games/action/show/id/1103632.html",
		}
		match2 := plusLigaMatch{
			homeTeam: "Aluron CMC Warta Zawiercie",
			awayTeam: "Steam Hemarpol Norwid Częstochowa",
			statsUrl: "plusliga-url/games/action/show/id/1103637.html",
		}
		plm := plusLigaMatches{
			newMatchKey(match1.homeTeam, match1.awayTeam): match1,
			newMatchKey(match2.homeTeam, match2.awayTeam): match2,
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
			domain.Match{
				ID:        6124,
				HomeName:  "Zawierce",
				AwayName:  "Norwid Czestochowa",
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
				StatsUrl:  match1.statsUrl,
			},
			domain.Match{
				ID:        6124,
				HomeName:  "Zawierce",
				AwayName:  "Norwid Czestochowa",
				StartTime: 1,
				Country:   "Poland",
				League:    "PlusLiga",
				StatsUrl:  match2.statsUrl,
			},
		}

		zipped, _ := plm.ZipWith(domainMatches)
		assert.Equal(t, expected, zipped)
	})
}
