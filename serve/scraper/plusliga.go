package scraper

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/p10r/pedro/serve/domain"
	"log"
	"net/http"
	"strings"
)

type PlusLiga struct {
	baseUrl string
}

type plusLigaMatch struct {
	homeTeam, awayTeam, statsUrl string
}

//nolint:unused
type matchKey string

//nolint:unused
func newMatchKey(homeName, awayName string) matchKey {
	return matchKey(homeName + "-" + awayName)
}

type plusLigaMatches []plusLigaMatch

var domainToPlusLigaMappings = map[string]string{
	"Barkom":             "Barkom Każany Lwów",
	"Bedzin":             "Nowak-Mosty MKS Będzin",
	"Belchatow":          "PGE GiEK SKRA Bełchatów",
	"Cuprum Gorzow":      "Cuprum Stilon Gorzów",
	"Gdansk":             "Trefl Gdańsk",
	"GKS Katowice":       "GKS Katowice",
	"Jastrzebski":        "Jastrzębski Węgiel",
	"Kedzierzyn-Kozle":   "ZAKSA Kędzierzyn-Koźle",
	"Lublin":             "BOGDANKA LUK Lublin",
	"Norwid Czestochowa": "Steam Hemarpol Norwid Częstochowa",
	"Olsztyn":            "Indykpol AZS Olsztyn",
	"Projekt Warszawa":   "PGE Projekt Warszawa",
	"Rzeszow":            "Asseco Resovia Rzeszów",
	"Stal Nysa":          "PSG Stal Nysa",
	"Slepsk Suwalki":     "Ślepsk Malow Suwałki",
	"Zawierce":           "Aluron CMC Warta Zawiercie",
}

//nolint:lll
func (m plusLigaMatches) ZipWith(dm domain.Matches) (zipped domain.Matches, notFound domain.Matches) {
	zipped = domain.Matches{}
	for _, d := range dm {
		plusLigaHome := domainToPlusLigaMappings[d.HomeName]
		plusLigaAway := domainToPlusLigaMappings[d.AwayName]

		for _, plMatch := range m {
			if plMatch.homeTeam == plusLigaHome && plMatch.awayTeam == plusLigaAway {
				d.StatsUrl = plMatch.statsUrl
				zipped = append(zipped, d)
			}
		}
	}

	return zipped, nil
}

func (pl *PlusLiga) ParseStats(res *http.Response) ([]plusLigaMatch, error) {
	var matches []plusLigaMatch

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return matches, err
	}

	doc.Find("section > .ajax-synced-games > .game-box").
		Each(func(i int, s *goquery.Selection) {
			// There should always exist one occurence of each inside .game-box
			homeSelection := s.Find("div > .game-team.left.gs").First().Text()
			awaySelection := s.Find("div > .game-team.right.gs").First().Text()
			statsSelection := s.Find("div > .game-more").First().Find("a")

			// Extracts /games/action/show/id/1103828.html
			statsUrl, statsExist := statsSelection.Attr("href")
			homeTeam := strings.TrimSpace(homeSelection)
			awayTeam := strings.TrimSpace(awaySelection)

			if awayTeam == "" || homeTeam == "" || !statsExist {
				// TODO: add some logging here
				return
			}

			match := plusLigaMatch{
				homeTeam: homeTeam,
				awayTeam: awayTeam,
				statsUrl: pl.baseUrl + statsUrl,
			}
			matches = append(matches, match)
		})

	return matches, nil
}

//nolint:unused
func (pl *PlusLiga) getAllMatchesPage() (*http.Response, error) {
	res, err := http.Get(pl.baseUrl + "/games.html")
	if err != nil {
		return nil, err
	}
	//defer res.Body.Close() TODO where?
	if res.StatusCode != 200 {
		log.Fatal("boom")
		// TODO logging
	}

	return res, nil
}

// printAllTeamsInLeague is intended to be run whenever there is a mismatch in teams.
// It is disabled by default.
//
//nolint:unused
//goland:noinspection GoUnusedFunction
func printAllTeamsInLeague(matches []plusLigaMatch) {
	teamsSet := map[string]string{}
	for _, match := range matches {
		teamsSet[match.homeTeam] = ""
		teamsSet[match.awayTeam] = ""
	}

	for key := range teamsSet {
		fmt.Println("", key)
	}
}
