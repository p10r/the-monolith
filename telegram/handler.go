package telegram

import (
	"gopkg.in/telebot.v3/middleware"
	"log"
	"pedro-go/db"
	"pedro-go/domain"
	"pedro-go/ra"
	"time"

	tele "gopkg.in/telebot.v3"
)

func Pedro(botToken, dsn string, allowedUserIds []int64) {
	now := func() time.Time { return time.Now() }

	conn := db.NewDB(dsn)
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	repo := db.NewSqliteArtistRepository(conn)

	m := db.NewEventMonitor(conn)
	r := domain.NewArtistRegistry(repo, ra.NewClient("https://ra.co/graphql"), m, now)

	bot, err := tele.NewBot(
		tele.Settings{
			Token:   botToken,
			Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
			Verbose: false,
		},
	)
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

	//bot.Use(middleware.Logger())
	bot.Handle("/follow", followArtist(r))
	bot.Handle("/list", listArtists(r))
	bot.Handle("/events", listEvents(r))
	bot.Start()
}
