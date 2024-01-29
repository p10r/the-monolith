package domain_test

import (
	"context"
	d "pedro-go/domain"
	"pedro-go/domain/expect"
	"pedro-go/ra"
	"sync/atomic"
	"testing"
	"time"
)

func TestMonitoring(t *testing.T) {
	ctx := context.Background()
	currentTs := atomic.Int64{}
	now := func() time.Time {
		t := time.UnixMilli(currentTs.Load())
		currentTs.Add(1000)
		return t
	}

	t.Run("records events", func(t *testing.T) {
		joe := d.UserID(444)
		eventId := "222"
		want := []d.MonitoringEvent{
			d.ArtistFollowed{
				ArtistSlug: "boysnoize",
				UserId:     joe,
				CreatedAt:  time.UnixMilli(0000),
			},
			d.NewEventForArtist{
				EventId:   "222",
				Slug:      "boysnoize",
				User:      joe,
				CreatedAt: time.UnixMilli(1000),
			},
		}

		registry, monitor := NewInMemoryArtistRegistry(
			map[d.RASlug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist: ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: ra.Events{
						{
							Id:         eventId,
							Title:      "Klubnacht 2",
							Date:       "2023-11-04T00:00:00.000",
							StartTime:  "2023-11-04T13:00:00.000",
							ContentUrl: "/events/1789025",
						},
					},
				},
			},
			now,
		)

		err := registry.Follow(ctx, "boysnoize", joe)
		expect.NoErr(t, err)

		_, err = registry.NewEventsForUser(ctx, joe)
		expect.NoErr(t, err)

		got, err := monitor.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})
}
