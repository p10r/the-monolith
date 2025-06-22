package telegram

import (
	"context"
	"fmt"
	"github.com/p10r/monolith/pedro/domain"
	"github.com/p10r/monolith/pkg/l"
	"gopkg.in/telebot.v3"
	"log/slog"
	"strconv"
	"time"
)

type Notifier struct {
	bot      *telebot.Bot
	registry *domain.ArtistRegistry
	users    []int64
	log      *slog.Logger
}

func NewNotifier(
	bot *telebot.Bot,
	registry *domain.ArtistRegistry,
	users []int64,
	log *slog.Logger,
) *Notifier {
	l := log.With(slog.String("adapter", "event_job"))

	return &Notifier{
		bot:      bot,
		registry: registry,
		users:    users,
		log:      l,
	}
}

func (n Notifier) StartEventNotifier() {
	ctx := context.Background()

	err := n.eventLookup()
	if err != nil {
		n.log.Error(l.Error("error when sending events to users", err))
	}

	duration := 12 * time.Hour
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	n.log.Info(fmt.Sprintf("event lookup is set to run every %v", duration.String()))

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		lookupErr := n.eventLookup()
		if lookupErr != nil {
			n.log.Error(l.Error("err when sending events", lookupErr))
		}
	}
}

func (n Notifier) eventLookup() error {
	ctx := context.Background()
	for _, id := range n.users {
		events, err := n.registry.NewEventsForUser(ctx, domain.UserID(id))
		if err != nil {
			return err
		}

		if len(events) == 0 {
			continue
		}

		n.log.Info(fmt.Sprintf("Sending %v", events))

		_, err = n.bot.Send(user(id), eventsMessage(events))
		if err != nil {
			return err
		}
	}
	return nil
}

type user int64

func (u user) Recipient() string {
	return strconv.FormatInt(int64(u), 10)
}
