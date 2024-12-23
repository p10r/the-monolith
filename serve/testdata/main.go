package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"strconv"
)

func main() {

	collyFun()
	//goQuery()
}

func goQuery() {
	// Request the HTML page.
	res, err := http.Get("https://www.plusliga.pl/games.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	occ := 0
	//<a class="btn btn-default btn-sm btn-more" href="/games/action
	doc.Find(".btn btn-default .btn-sm btn-more").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the title
		//title := s.Find("a").Text()
		occ++
	})
	fmt.Printf(strconv.Itoa(occ))
}

type plusLigaMatch struct {
	// Fetched from data-game-id, e.g. 1103467
	gameId string
	// e.g. https://www.plusliga.pl/games/action/show/id/1103467.html
	statsUrl string
}

func collyFun() {
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"))

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML("section > .ajax-synced-games ", func(e *colly.HTMLElement) {

		//println("yeye")

		//1103467
		///games/action/show/id/1103467.html
		gameId := e.Attr("data-game-id")
		//homeTeam := e.ChildAttr(".data-synced-games-class", "teamAGSStatus")
		//awayTeam := e.Attr("teamBGSStatus")

		fmt.Println(gameId)
		//fmt.Println(homeTeam)
		//fmt.Println(awayTeam)

		//fmt.Println("yeet")
	})

	c.Visit("https://www.plusliga.pl/games.html")
}
