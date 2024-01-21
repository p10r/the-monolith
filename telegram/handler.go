package telegram

import (
	"context"
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

	bot.Use(middleware.Whitelist(allowedUserIds...))

	log.Print("Started Pedro")

	n := &Notifier{
		bot:      bot,
		registry: r,
		users:    allowedUserIds,
	}

	n.NotifyUsers()

	bot.Use(middleware.Logger())
	bot.Handle("/follow", followArtist(r))
	bot.Handle("/list", listArtists(r))
	bot.Handle("/events", listEvents(r))
	bot.Start()
}

func listEvents(r *domain.ArtistRegistry) func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background() //TODO check if telebot can provide context

		artists, err := r.ArtistsFor(ctx, domain.UserID(c.Sender().ID))
		eventsFor, err := r.AllEventsForArtist(ctx, artists[0])
		if err != nil {
			log.Print(err)
			return c.Send("There was an error!")
		}

		return c.Send(fmt.Sprintf("Event:\n%v", eventsFor))
	}
}

func listArtists(r *domain.ArtistRegistry) func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background() //TODO check if telebot can provide context

		artists, err := r.ArtistsFor(ctx, domain.UserID(c.Sender().ID))
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
	}
}

func followArtist(r *domain.ArtistRegistry) func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background() //TODO check if telebot can provide context

		tags := c.Args()
		slug, err := domain.NewSlug(tags[0])
		if err != nil {
			log.Print(err)
			return c.Send("Could not parse artist, make sure to send it as follows https://ra.co/dj/yourartist")
		}
		userId := domain.UserID(c.Sender().ID)

		log.Printf("%v started following %v", userId, slug)
		err = r.Follow(ctx, slug, userId)
		if err != nil {
			log.Print(err)
			return c.Send("There was an error!")
		}

		return c.Send("Done!")
	}
}
