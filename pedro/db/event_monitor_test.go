package db

import (
	"context"
	d "github.com/p10r/pedro/pedro/domain"
	"github.com/p10r/pedro/pedro/domain/expect"
	"testing"
)

func TestEvents(t *testing.T) {
	ctx := context.Background()

	t.Run("records events", func(t *testing.T) {
		sqlite := MustOpenDB(t)
		monitor := NewEventMonitor(sqlite)

		want := []d.MonitoringEvent{
			d.ArtistFollowed{ArtistSlug: "sinamin", UserId: 666},
			d.NewEventForArtist{EventId: "123", Slug: "sinamin", User: 1},
		}

		monitor.Monitor(ctx, d.ArtistFollowed{ArtistSlug: "sinamin", UserId: 666})
		monitor.Monitor(ctx, d.NewEventForArtist{EventId: "123", Slug: "sinamin", User: 1})

		got, err := monitor.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})
}
