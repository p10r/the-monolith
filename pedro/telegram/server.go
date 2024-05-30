package telegram

import (
	"github.com/p10r/pedro/pedro/db"
	"github.com/p10r/pedro/pedro/domain"
	"github.com/p10r/pedro/pedro/ra"
	"github.com/p10r/pedro/pkg/sqlite"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"time"
)

// NewPedro wires Pedro App together.
// Expects an already opened connection.
func NewPedro(
	conn *sqlite.DB,
	botToken string,
	allowedUserIds []int64,
) *telebot.Bot {
	now := func() time.Time { return time.Now() }

	repo := db.NewSqliteArtistRepository(conn)

	m := db.NewEventMonitor(conn)
	artistRegistry := domain.NewArtistRegistry(repo, ra.NewClient("https://ra.co"), m, now)

	bot, err := telebot.NewBot(
		telebot.Settings{
			Token:   botToken,
			Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
			Verbose: false,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	bot.Use(middleware.Whitelist(allowedUserIds...))

	log.Print("Started Pedro")

	n := &Notifier{
		bot:      bot,
		registry: artistRegistry,
		users:    allowedUserIds,
	}

	go n.StartEventNotifier()

	sender := TelebotSender{}

	//bot.Use(middleware.Logger())
	bot.Handle("/follow", followArtist(artistRegistry, sender))
	bot.Handle("/artists", listArtists(artistRegistry, sender))
	bot.Handle("/events", listEvents(artistRegistry, sender))

	return bot
}
