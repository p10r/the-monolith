package statistics

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/p10r/pedro/serve/domain"
	"net/http"
	"regexp"
	"strings"
)

// Every team links to a team page on the SuperLega page.
// Instead of parsing the string, we use the id of the match page.
// E.g.: Milano: http://www.legavolley.it/team/6551
var domainToSuperLegaIds = map[string]string{
	"Milano":          "6551",
	"Cisterna":        "6552",
	"Lube Civitanova": "6553",
	"Piacenza":        "6554",
	"Taranto":         "6555",
	"Trentino":        "6556",
	"Monza":           "6557",
	"Padova":          "6558",
	"Verona":          "6559",
	"Perugia":         "6560",
	"Modena":          "6561",
	"Grottazzolina":   "6592",
}

type superLegaMatch struct {
	homeTeamId, awayTeamId, statsUrl string
}

type superLegaMatches map[matchKey]superLegaMatch

func (m superLegaMatches) ZipWith(
	domainMatches domain.Matches,
) (zipped domain.Matches, notFound domain.Matches) {
	zipped = domain.Matches{}
	notFound = domain.Matches{}

	for _, dm := range domainMatches {
		superLegaHome := domainToSuperLegaIds[dm.HomeName]
		superLegaAway := domainToSuperLegaIds[dm.AwayName]
		key := newMatchKey(superLegaHome, superLegaAway)

		slm, ok := m[key]
		if !ok {
			notFound = append(notFound, dm)
			continue
		}

		dm.StatsUrl = slm.statsUrl
		zipped = append(zipped, dm)
	}

	return zipped, nil
}

type superLegaScraper struct {
	baseUrl string
	client  *http.Client
}

func newSuperLegaScraper(baseUrl string, client *http.Client) *superLegaScraper {

	return &superLegaScraper{baseUrl: strings.TrimSuffix(baseUrl, "/"), client: client}
}

func (scraper *superLegaScraper) GetStatsFor(dm domain.Matches) (
	matched domain.Matches,
	notFound domain.Matches,
	err error,
) {
	page, err := scraper.getAllMatchesPage()
	if err != nil {
		return nil, nil, fmt.Errorf("could not fetch superLegaScraper match page: %w", err)
	}

	plMatches, err := scraper.parseStats(page)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse superLegaScraper match page: %w", err)
	}

	zipped, notFound := plMatches.ZipWith(dm)
	return zipped, notFound, nil
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

				matches[newMatchKey(homeTeamId, awayTeamId)] = superLegaMatch{
					homeTeamId: homeTeamId,
					awayTeamId: awayTeamId,
					statsUrl:   "http://www.legavolley.it/match/" + statId,
				}
			})
		})
	return matches, nil
}
