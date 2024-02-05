package telegram

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"pedro-go/domain"
	"strings"
)

func eventsMessage(events domain.Events) string {
	if len(events) == 0 {
		return fmt.Sprintln("There are no events in the near future.")
	}

	var lines []string
	for _, e := range events {
		layout := "02.01 15:04"
		line := fmt.Sprintf("%v - %v@%v", e.StartTime.Format(layout), e.Artist, e.Venue)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func listEvents(r *domain.ArtistRegistry) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		ctx := context.Background() //TODO check if telebot can provide context

		events, err := r.NewEventsForUser(ctx, domain.UserID(c.Sender().ID))
		if err != nil {
			log.Print(err)
			return c.Send(genericErrMsg("/events", err))
		}

		return c.Send(eventsMessage(events))
	}
}
