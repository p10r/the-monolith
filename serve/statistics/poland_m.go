package statistics

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/p10r/monolith/serve/domain"
	"net/http"
	"strings"
)

type plusLigaMatches []plusLigaMatch

func (m plusLigaMatches) toDomain() domain.StatSheets {
	statSheets := domain.StatSheets{}
	for _, match := range m {
		statSheets = append(statSheets, match.toDomain())
	}
	return statSheets
}

type plusLigaMatch struct {
	homeTeam plusligaTeamId
	awayTeam plusligaTeamId
	statsUrl string
}

func (m plusLigaMatch) toDomain() domain.StatSheet {
	return domain.StatSheet{
		Home: m.homeTeam.toDomain(),
		Away: m.awayTeam.toDomain(),
		Url:  m.statsUrl,
	}
}

type plusligaTeamId string

func (t plusligaTeamId) toDomain() string {
	return plusLigaMappings[t]
}

var plusLigaMappings = map[plusligaTeamId]string{
	"Barkom Każany Lwów":                "Barkom",
	"Nowak-Mosty MKS Będzin":            "Bedzin",
	"PGE GiEK SKRA Bełchatów":           "Belchatow",
	"Cuprum Stilon Gorzów":              "Cuprum Gorzow",
	"Trefl Gdańsk":                      "Gdansk",
	"GKS Katowice":                      "GKS Katowice",
	"Jastrzębski Węgiel":                "Jastrzebski",
	"ZAKSA Kędzierzyn-Koźle":            "Kedzierzyn-Kozle",
	"BOGDANKA LUK Lublin":               "Lublin",
	"Steam Hemarpol Norwid Częstochowa": "Norwid Czestochowa",
	"Indykpol AZS Olsztyn":              "Olsztyn",
	"PGE Projekt Warszawa":              "Projekt Warszawa",
	"Asseco Resovia Rzeszów":            "Rzeszow",
	"PSG Stal Nysa":                     "Stal Nysa",
	"Ślepsk Malow Suwałki":              "Slepsk Suwalki",
	"Aluron CMC Warta Zawiercie":        "Zawierce",
}

type plusLigaScraper struct {
	baseUrl string
	client  *http.Client
}

func newPlusLiga(baseUrl string, client *http.Client) *plusLigaScraper {
	return &plusLigaScraper{baseUrl: strings.TrimSuffix(baseUrl, "/"), client: client}
}

func (scraper *plusLigaScraper) GetStats() (plusLigaMatches, error) {
	page, err := scraper.getAllMatchesPage()
	if err != nil {
		return nil, fmt.Errorf("could not fetch plusLigaScraper match page: %w", err)
	}

	plMatches, err := scraper.parseStats(page)
	if err != nil {
		return nil, fmt.Errorf("could not parse plusLigaScraper match page: %w", err)
	}

	return plMatches, nil
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

			match := plusLigaMatch{
				homeTeam: plusligaTeamId(homeTeam),
				awayTeam: plusligaTeamId(awayTeam),
				statsUrl: "https://www.plusliga.pl" + statsUrl,
			}
			matches = append(matches, match)
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
	teamsSet := map[plusligaTeamId]string{}
	for _, match := range matches {
		teamsSet[match.homeTeam] = ""
		teamsSet[match.awayTeam] = ""
	}

	for key := range teamsSet {
		fmt.Println("", key)
	}
}
