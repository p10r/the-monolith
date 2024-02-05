package telegram

import (
	"pedro-go/domain"
	"testing"
	"time"
)

func TestMessages(t *testing.T) {

	t.Run("creates valid events message", func(t *testing.T) {
		eventsMessage(
			domain.Events{
				domain.Event{
					Id:         0,
					Title:      "Event 1",
					Artist:     "Boys Noize",
					Venue:      "Gretchen",
					StartTime:  time.Date(2024, 1, 12, 23, 0, 0, 0, time.UTC),
					ContentUrl: "",
				},
			})
	})

}
