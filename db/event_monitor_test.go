package db

import (
	"context"
	d "pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestEvents(t *testing.T) {
	ctx := context.Background()

	t.Run("records events", func(t *testing.T) {
		sqlite := MustOpenDB(t)
		monitor := NewEventMonitor(sqlite)

		want := []d.MonitoringEvent{
			d.ArtistFollowedEvent{ArtistSlug: "sinamin", UserId: 666},
			d.NewEventForArtistEvent{EventId: "123", Slug: "sinamin", Users: []d.UserID{1, 3}},
		}

		monitor.Monitor(ctx, d.ArtistFollowedEvent{ArtistSlug: "sinamin", UserId: 666})
		monitor.Monitor(ctx, d.NewEventForArtistEvent{EventId: "123", Slug: "sinamin", Users: []d.UserID{1, 3}})

		got, err := monitor.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})
}
