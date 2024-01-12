package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"pedro-go/db"
	"pedro-go/domain"
	"pedro-go/ra"
	"strings"
	"time"

	tele "gopkg.in/telebot.v3"
)

func Pedro(botToken, dsn string, allowedUserIds []int64) {
	repo, err := db.NewGormArtistRepository(dsn)
	if err != nil {
		log.Fatalf("Cannot connect to db %v", err)
	}

	r := domain.NewArtistRegistry(repo, ra.NewClient("https://ra.co/graphql"))

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

	bot.Use(middleware.Logger())
	bot.Use(middleware.Whitelist(allowedUserIds...))

	bot.Handle("/follow", func(c tele.Context) error {
		tags := c.Args()
		slug, err := ra.NewSlug(tags[0])
		if err != nil {
			log.Print(err)
			return c.Send("Could not parse artist, make sure to send it as follows https://ra.co/dj/yourartist")
		}
		userId := domain.UserId(c.Sender().ID)

		log.Printf("%v started following %v", userId, slug)
		err = r.Follow(slug, userId)
		if err != nil {
			log.Print(err)
			return c.Send("There was an error!")
		}

		return c.Send("Done!")
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

	bot.Handle("/events", func(c tele.Context) error {
		artists, err := r.ArtistsFor(domain.UserId(c.Sender().ID))
		eventsFor, err := r.EventsFor(artists[0])
		if err != nil {
			log.Print(err)
			return c.Send("There was an error!")
		}

		return c.Send(fmt.Sprintf("Event:\n%v", eventsFor))
	})

	log.Print("Started Pedro")
	bot.Start()
}
