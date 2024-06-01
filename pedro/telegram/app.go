package telegram

import (
	"github.com/p10r/pedro/pedro/db"
	"github.com/p10r/pedro/pedro/domain"
	"github.com/p10r/pedro/pedro/ra"
	"github.com/p10r/pedro/pkg/sqlite"
	"gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log/slog"
	"time"
)

// NewPedro wires Pedro App together.
// Expects an already opened connection.
func NewPedro(
	conn *sqlite.DB,
	botToken string,
	allowedUserIds []int64,
	logHandler slog.Handler,
) *telebot.Bot {
	log := slog.New(logHandler).With(slog.String("app", "pedro"))

	artistRegistry := domain.NewArtistRegistry(
		db.NewSqliteArtistRepository(conn),
		ra.NewClient("https://ra.co", log),
		func() time.Time { return time.Now() },
		log,
	)

	bot, err := telebot.NewBot(
		telebot.Settings{
			Token:   botToken,
			Poller:  &telebot.LongPoller{Timeout: 10 * time.Second},
			Verbose: false,
		},
	)
	if err != nil {
		log.Error("cannot create telegram bot", slog.Any("error", err))
	}

	bot.Use(middleware.Whitelist(allowedUserIds...))

	n := NewNotifier(bot, artistRegistry, allowedUserIds, log)
	go n.StartEventNotifier()

	sender := NewTelegramSender(log)

	l := log.With(slog.String("adapter", "telegram_in"))
	//bot.Use(middleware.Logger()) //TODO check if slog can be added here
	bot.Handle("/follow", followArtist(artistRegistry, sender, l))
	bot.Handle("/artists", listArtists(artistRegistry, sender, l))
	bot.Handle("/events", listEvents(artistRegistry, sender, l))

	log.Info("started pedro")

	return bot
}
