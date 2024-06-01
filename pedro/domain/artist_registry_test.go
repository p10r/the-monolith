package domain_test

import (
	"context"
	"github.com/p10r/pedro/pedro/db"
	"github.com/p10r/pedro/pedro/domain"
	"github.com/p10r/pedro/pedro/domain/expect"
	"github.com/p10r/pedro/pedro/ra"
	"github.com/p10r/pedro/pkg/l"
	"testing"
	"time"
)

func NewInMemoryArtistRegistry(
	t *testing.T,
	raArtists ra.ArtistStore,
	now func() time.Time,
) *domain.ArtistRegistry {
	repo := db.NewInMemoryArtistRepository()
	raClient := ra.NewInMemoryClient(t, raArtists)
	log := l.NewTextLogger()
	return domain.NewArtistRegistry(repo, raClient, now, log)
}

func TestArtistRegistry(t *testing.T) {
	ctx := context.Background()
	now := func() time.Time {
		return time.Now()
	}

	t.Run("lists all artists", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{},
				},
			},
			now,
		)

		err := registry.Follow(ctx, "boysnoize", domain.UserID(1))
		expect.NoErr(t, err)

		err = registry.Follow(ctx, "sinamin", domain.UserID(1))
		expect.NoErr(t, err)

		got, err := registry.All(ctx)
		expect.NoErr(t, err)
		want := domain.Artists{
			{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    domain.UserIDs{domain.UserID(1)},
				TrackedEvents: domain.EventIDs{},
			},
			{
				ID:            2,
				RAID:          "222",
				RASlug:        "sinamin",
				Name:          "Sinamin",
				FollowedBy:    domain.UserIDs{domain.UserID(1)},
				TrackedEvents: domain.EventIDs{},
			},
		}

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("follows an artist from resident advisor", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"daftpunk": {
					Artist:     ra.Artist{RAID: "111", Name: "Daft Punk"},
					EventsData: []ra.Event{},
				},
			},
			now,
		)

		want := []domain.Artist{
			{
				ID:            1,
				RAID:          "111",
				RASlug:        "daftpunk",
				Name:          "Daft Punk",
				FollowedBy:    domain.UserIDs{domain.UserID(1)},
				TrackedEvents: domain.EventIDs{},
			}}
		err := registry.Follow(ctx, "daftpunk", domain.UserID(1))

		expect.NoErr(t, err)
		all, err := registry.All(ctx)
		expect.NoErr(t, err)
		expect.DeepEqual(t, all, want)
	})

	t.Run("doesn't add artist if already added", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
			},
			now,
		)

		want := domain.Artists{
			{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    domain.UserIDs{domain.UserID(1)},
				TrackedEvents: domain.EventIDs{},
			},
		}
		err := registry.Follow(ctx, "boysnoize", domain.UserID(1))

		expect.NoErr(t, err)
		all, err := registry.All(ctx)
		expect.NoErr(t, err)
		expect.DeepEqual(t, all, want)
	})

	t.Run("returns error if artist can't be found on RA", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{},
			now,
		)

		err := registry.Follow(ctx, "unknown", domain.UserID(1))

		expect.Err(t, err)
		expect.Equal(t, err.Error(), domain.ErrNotFoundOnRA.Error())
	})

	t.Run("follows new artist as user", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{},
				},
			},
			now,
		)

		err := registry.Follow(ctx, "boysnoize", domain.UserID(1))
		expect.NoErr(t, err)

		err = registry.Follow(ctx, "sinamin", domain.UserID(2))
		expect.NoErr(t, err)

		got, err := registry.ArtistsFor(ctx, domain.UserID(1))
		want := domain.Artists{
			domain.Artist{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    domain.UserIDs{domain.UserID(1)},
				TrackedEvents: domain.EventIDs{},
			},
		}

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("ignores follow if already following", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{},
				},
			},
			now,
		)
		err := registry.Follow(ctx, "boysnoize", domain.UserID(1))
		expect.NoErr(t, err)

		err = registry.Follow(ctx, "boysnoize", domain.UserID(1))
		expect.NoErr(t, err)

		got, err := registry.ArtistsFor(ctx, domain.UserID(1))

		expect.NoErr(t, err)
		expect.True(t, len(got) == 1)
		expect.DeepEqual(t, got[0].FollowedBy, domain.UserIDs{domain.UserID(1)})
	})

	t.Run("follows existing artist", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{},
				},
			},
			now,
		)
		err := registry.Follow(ctx, "boysnoize", domain.UserID(1))
		expect.NoErr(t, err)

		err = registry.Follow(ctx, "boysnoize", domain.UserID(2))
		expect.NoErr(t, err)

		got, err := registry.ArtistsFor(ctx, domain.UserID(2))

		expect.NoErr(t, err)
		expect.Equal(t, len(got), 1)
		expect.DeepEqual(t, got[0].FollowedBy, domain.UserIDs{domain.UserID(1), domain.UserID(2)})
	})

	t.Run("fetches all events for artist in the next month", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist: ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{
						{
							Id:         "1",
							Title:      "Klubnacht",
							StartTime:  "2023-11-04T13:00:00.000",
							ContentUrl: "/events/1789025",
						},
						{
							Id:         "2",
							Title:      "Klubnacht 2",
							StartTime:  "2023-11-04T13:00:00.000",
							ContentUrl: "/events/1789025",
						},
					},
				},
			},
			now,
		)

		events, err := registry.EventsForArtist(ctx, domain.Artist{
			ID:         1,
			RAID:       "222",
			RASlug:     "sinamin",
			Name:       "Sinamin",
			FollowedBy: nil,
		})

		expect.NoErr(t, err)
		expect.Equal(t, len(events), 2)
	})

	t.Run("reports only new events to user", func(t *testing.T) {
		events := ra.Events{
			{
				Id:         "3",
				Title:      "Kater Blau Night",
				StartTime:  "2023-11-04T13:00:00.000",
				ContentUrl: "/events/3",
				Venue: ra.Venue{
					Area: ra.Area{Name: "Berlin"},
					Name: "RSO",
				},
			},
			{
				Id:         "1",
				Title:      "Klubnacht",
				StartTime:  "2023-11-04T13:00:00.000",
				ContentUrl: "/events/1789025",
				Venue: ra.Venue{
					Area: ra.Area{Name: "Berlin"},
					Name: "RSO",
				},
			},
			{
				Id:         "2",
				Title:      "Klubnacht 2",
				StartTime:  "2023-11-04T13:00:00.000",
				ContentUrl: "/events/1789025",
				Venue: ra.Venue{
					Area: ra.Area{Name: "Berlin"},
					Name: "RSO",
				},
			},
		}
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{events[0]},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{events[1], events[2]},
				},
			},
			now,
		)

		joe := domain.UserID(1)

		err := registry.Follow(ctx, "boysnoize", joe)
		expect.NoErr(t, err)

		err = registry.Follow(ctx, "sinamin", joe)
		expect.NoErr(t, err)

		eventsForUser, err := registry.NewEventsForUser(ctx, joe)
		expect.NoErr(t, err)

		nov4th1pm := time.Date(2023, 11, 4, 13, 0, 0, 0, time.UTC)
		expect.DeepEqual(t, eventsForUser, domain.Events{
			{
				Id:         domain.EventID(3),
				Title:      "Kater Blau Night",
				Artist:     "Boys Noize",
				Venue:      "RSO",
				City:       "Berlin",
				StartTime:  nov4th1pm,
				ContentUrl: "/events/3",
			},
			{
				Id:         domain.EventID(1),
				Title:      "Klubnacht",
				Artist:     "Sinamin",
				Venue:      "RSO",
				City:       "Berlin",
				StartTime:  nov4th1pm,
				ContentUrl: "/events/1789025",
			},
			{
				Id:         domain.EventID(2),
				Title:      "Klubnacht 2",
				Artist:     "Sinamin",
				Venue:      "RSO",
				City:       "Berlin",
				StartTime:  nov4th1pm,
				ContentUrl: "/events/1789025",
			},
		},
		)

		newlyFetched, err := registry.NewEventsForUser(ctx, joe)

		expect.NoErr(t, err)
		expect.Len(t, newlyFetched, 0)
	})

	t.Run("allEvents always fetches all events for user", func(t *testing.T) {
		events := ra.Events{
			{
				Id:         "3",
				Title:      "Kater Blau Night",
				StartTime:  "2023-11-04T13:00:00.000",
				ContentUrl: "/events/3",
				Venue: ra.Venue{
					Area: ra.Area{Name: "Berlin"},
					Name: "RSO",
				},
			},
		}
		registry := NewInMemoryArtistRegistry(t,
			ra.ArtistStore{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{events[0]},
				},
			},
			now,
		)

		joe := domain.UserID(1)

		err := registry.Follow(ctx, "boysnoize", joe)
		expect.NoErr(t, err)

		_, err = registry.AllEventsForUser(ctx, joe)
		expect.NoErr(t, err)

		got, err := registry.AllEventsForUser(ctx, joe)
		expect.NoErr(t, err)

		log := l.NewTextLogger()
		expect.DeepEqual(t, got, events.ToDomainEvents("Boys Noize", log))
	})
}
