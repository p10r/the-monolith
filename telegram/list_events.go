package telegram

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"pedro-go/domain"
	"strings"
)

func newEventsMessage(events domain.Events) string {
	var lines []string
	for _, e := range events {
		line := fmt.Sprintf("%v - %v@%v", e.Date, e.Artist, e.Venue)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func noUpcomingEventsMsg(events domain.Events) string {
	return fmt.Sprintln("There are no events in the near future.", events)
}

func listEvents(r *domain.ArtistRegistry) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		ctx := context.Background() //TODO check if telebot can provide context

		events, err := r.NewEventsForUser(ctx, domain.UserID(c.Sender().ID))
		if err != nil {
			log.Print(err)
			return c.Send(genericErrMsg("/events", err))
		}

		if len(events) == 0 {
			return c.Send(noUpcomingEventsMsg(events))
		}

		return c.Send(newEventsMessage(events))
	}
}
