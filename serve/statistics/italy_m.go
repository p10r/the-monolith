package statistics

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/p10r/pedro/serve/domain"
	"net/http"
	"regexp"
	"strings"
)

type superLegaTeamId string

func (i superLegaTeamId) toDomain() string {
	return superLegaToDomainMappings[i]
}

// Every team links to a team page on the SuperLega page.
// Instead of parsing the string, we use the id of the match page.
// E.g.: Milano: http://www.legavolley.it/team/6551
var superLegaToDomainMappings = map[superLegaTeamId]string{
	"6551": "Milano",
	"6552": "Cisterna",
	"6553": "Lube Civitanova",
	"6554": "Piacenza",
	"6555": "Taranto",
	"6556": "Trentino",
	"6557": "Monza",
	"6558": "Padova",
	"6559": "Verona",
	"6560": "Perugia",
	"6561": "Modena",
	"6592": "Grottazzolina",
}

type superLegaMatch struct {
	homeTeamId superLegaTeamId
	awayTeamId superLegaTeamId
	statsUrl   string
}

func (m superLegaMatch) toDomain() domain.StatSheet {
	return domain.StatSheet{
		League: "italy: superlega",
		Home:   m.homeTeamId.toDomain(),
		Away:   m.awayTeamId.toDomain(),
		Url:    m.statsUrl,
	}
}

type superLegaMatches []superLegaMatch

func (m superLegaMatches) toDomain() domain.StatSheets {
	statSheets := domain.StatSheets{}
	for _, match := range m {
		statSheets = append(statSheets, match.toDomain())
	}
	return statSheets
}

type superLegaScraper struct {
	baseUrl string
	client  *http.Client
}

func newSuperLegaScraper(baseUrl string, client *http.Client) *superLegaScraper {
	return &superLegaScraper{baseUrl: strings.TrimSuffix(baseUrl, "/"), client: client}
}

func (scraper *superLegaScraper) GetStats() (superLegaMatches, error) {
	page, err := scraper.getAllMatchesPage()
	if err != nil {
		return nil, fmt.Errorf("could not fetch superLegaScraper match page: %w", err)
	}

	matches, err := scraper.parseStats(page)
	if err != nil {
		return nil, fmt.Errorf("could not parse superLegaScraper match page: %w", err)
	}

	return matches, nil
}

func (scraper *superLegaScraper) getAllMatchesPage() (*http.Response, error) {
	//https://www.legavolley.it/calendario/?lang=en
	res, err := scraper.client.Get(scraper.baseUrl + "/calendario/?lang=en")
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got %d when fetching page. err: %w", res.StatusCode, err)
	}

	return res, nil
}

var superLegaUrlRegex = regexp.MustCompile(`\d+`)

//nolint:lll
func (scraper *superLegaScraper) parseStats(res *http.Response) (matches superLegaMatches, err error) {
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return matches, err
	}

	matches = superLegaMatches{}
	doc.Find("table#GareGiornata").
		Each(func(i int, selection *goquery.Selection) {
			// Drop wrapper table
			if i == 0 {
				return
			}

			// Each match day
			selection.Find(".tab-gara").Each(func(i2 int, tr *goquery.Selection) {
				// Each <tr> has multiple tds
				//        <td> Home Team
				//        <td> Results + Stats URL
				//        <td> Away Team
				//        <td> Refs
				//        <td> Location
				//        <td> Location
				//        <td> Generic Streaming Link
				//
				// The ones we're interested in are the first three:
				//  - Home Team
				//  - Results + Stats URL
				//  - Away Team
				//
				// The information that's important is stored inside the 'onClick' action
				onClickContent := []string{}
				for _, n := range tr.ChildrenFiltered("td").Nodes {
					for _, attr := range n.Attr {
						if attr.Key == "onclick" {
							onClickContent = append(onClickContent, attr.Val)
						}
					}
				}

				if len(onClickContent) < 3 {
					return
				}

				homeTeamId := superLegaUrlRegex.FindString(onClickContent[0])
				awayTeamId := superLegaUrlRegex.FindString(onClickContent[2])
				statId := superLegaUrlRegex.FindString(onClickContent[1])

				match := superLegaMatch{
					homeTeamId: superLegaTeamId(homeTeamId),
					awayTeamId: superLegaTeamId(awayTeamId),
					statsUrl:   "http://www.legavolley.it/match/" + statId,
				}
				matches = append(matches, match)
			})
		})
	return matches, nil
}
