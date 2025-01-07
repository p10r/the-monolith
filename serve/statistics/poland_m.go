package statistics

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/p10r/pedro/serve/domain"
	"net/http"
	"strings"
)

type plusLigaMatch struct {
	homeTeam, awayTeam, statsUrl string
}

type plusLigaMatches map[matchKey]plusLigaMatch

type matchKey string

func newMatchKey(homeName, awayName string) matchKey {
	return matchKey(homeName + "-" + awayName)
}

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

func (m plusLigaMatches) ZipWith(
	dm domain.Matches,
) (zipped domain.Matches, notFound domain.Matches) {
	zipped = domain.Matches{}
	notFound = domain.Matches{}

	for _, d := range dm {
		plusLigaHome := domainToPlusLigaMappings[d.HomeName]
		plusLigaAway := domainToPlusLigaMappings[d.AwayName]
		key := newMatchKey(plusLigaHome, plusLigaAway)

		plMatch, ok := m[key]
		if !ok {
			notFound = append(notFound, d)
			continue
		}

		d.StatsUrl = plMatch.statsUrl
		zipped = append(zipped, d)
	}

	return zipped, nil
}

type plusLigaScraper struct {
	baseUrl string
	client  *http.Client
}

func newPlusLiga(baseUrl string, client *http.Client) *plusLigaScraper {
	return &plusLigaScraper{baseUrl: strings.TrimSuffix(baseUrl, "/"), client: client}
}

func (scraper *plusLigaScraper) GetAllStats() (statUrls []string, err error) {
	page, err := scraper.getAllMatchesPage()
	if err != nil {
		return []string{}, fmt.Errorf("could not fetch plusLigaScraper match page: %w", err)
	}

	plMatches, err := scraper.parseStats(page)
	if err != nil {
		return []string{}, fmt.Errorf("could not parse plusLigaScraper match page: %w", err)
	}

	urls := []string{}
	for _, match := range plMatches {
		urls = append(urls, match.statsUrl)
	}
	return urls, nil
}

func (scraper *plusLigaScraper) GetStatsFor(dm domain.Matches) (
	matched domain.Matches,
	notFound domain.Matches,
	err error,
) {
	page, err := scraper.getAllMatchesPage()
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch plusLigaScraper match page: %w", err)
	}

	plMatches, err := scraper.parseStats(page)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse plusLigaScraper match page: %w", err)
	}

	zipped, notFound := plMatches.ZipWith(dm)
	return zipped, notFound, nil
}

func (scraper *plusLigaScraper) getAllMatchesPage() (*http.Response, error) {
	res, err := scraper.client.Get(scraper.baseUrl + "/games.html")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got %d when fetching page. err: %w", res.StatusCode, err)
	}

	return res, nil
}

//nolint:lll
func (scraper *plusLigaScraper) parseStats(res *http.Response) (matches plusLigaMatches, err error) {
	defer res.Body.Close()
	matches = plusLigaMatches{}

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
				return
			}

			matches[newMatchKey(homeTeam, awayTeam)] = plusLigaMatch{
				homeTeam: homeTeam,
				awayTeam: awayTeam,
				statsUrl: scraper.baseUrl + statsUrl,
			}
		})

	//printAllTeamsInLeague(matches)

	return matches, nil
}

// printAllTeamsInLeague is intended to be run whenever there is a mismatch in teams.
// It is disabled by default.
//
//nolint:unused
//goland:noinspection GoUnusedFunction
func printAllTeamsInLeague(matches plusLigaMatches) {
	teamsSet := map[string]string{}
	for _, match := range matches {
		teamsSet[match.homeTeam] = ""
		teamsSet[match.awayTeam] = ""
	}

	for key := range teamsSet {
		fmt.Println("", key)
	}
}
