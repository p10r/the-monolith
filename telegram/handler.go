package main

import (
	"fmt"
	"log"
	"os"
	"pedro-go/db"
	"pedro-go/domain"
	"pedro-go/ra"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	Pedro(token)
}

func Pedro(botToken string) {
	r := domain.NewArtistRegistry(
		db.NewInMemoryArtistRepository(),
		ra.NewClient("https://ra.co/graphql"),
	)

	pref := tele.Settings{
		Token:   botToken,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: false,
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/follow", func(c tele.Context) error {
		tags := c.Args()

		_, err = r.Follow(ra.Slug(tags[0]), domain.UserId(c.Sender().ID))
		if err != nil {
			log.Print(err)
			return c.Send("There was an error!")
		}

		return c.Send("Hello!")
	})

	bot.Handle("/list", func(c tele.Context) error {
		artists, err := r.ArtistsFor(domain.UserId(c.Sender().ID))
		if err != nil {
			log.Print(err)
			return c.Send("There was an error!")
		}

		var res []string
		for _, artist := range artists {
			res = append(res, "- "+artist.Name)
		}

		if len(res) == 0 {
			return c.Send("You're not following anyone yet.")
		}

		return c.Send(fmt.Sprintf("You're following:\n%v", strings.Join(res, "\n")))
	})

	log.Print("Started Pedro")
	bot.Start()
}
