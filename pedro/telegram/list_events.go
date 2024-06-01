package telegram

import (
	"context"
	"fmt"
	"github.com/p10r/pedro/pedro/domain"
	"gopkg.in/telebot.v3"
	"log/slog"
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

func listEvents(
	r *domain.ArtistRegistry,
	s Sender,
	log *slog.Logger,
) func(c telebot.Context) error {
	return func(c telebot.Context) error {
		ctx := context.Background() //TODO check if telebot can provide context

		id := c.Sender().ID
		events, err := r.AllEventsForUser(ctx, domain.UserID(id))
		if err != nil {
			log.Error(
				fmt.Sprintf("%v has no upcoming events", id),
				slog.Any("error", err),
			)
			return s.Send(c, genericErrMsg("/events", err))
		}

		return s.Send(c, eventsMessage(events))
	}
}
