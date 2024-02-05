package telegram

import (
	"context"
	"gopkg.in/telebot.v3"
	"log"
	"pedro-go/domain"
	"time"
)

type Notifier struct {
	bot      *telebot.Bot
	registry *domain.ArtistRegistry
	users    []int64
}

func (n Notifier) StartEventNotifier() {
	ctx := context.Background()
	if ctx == nil {
		log.Fatal("Event lookup: ctx setup failed")
	}
	duration := 12 * time.Hour
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	log.Printf("event lookup is set to run every %v", duration.String())

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}

		err := n.eventLookup()
		if err != nil {
			log.Printf("Event lookup: Error when sending events to users. err: %v", err)
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

		log.Printf("Event lookup: Sending %v\n", events)
		_, err = n.bot.Send(user("530586914"), eventsMessage(events))
		if err != nil {
			return err
		}
	}
	return nil
}

type user string

func (u user) Recipient() string {
	return string(u)
}
