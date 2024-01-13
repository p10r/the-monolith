package domain_test

import (
	"pedro-go/db"
	. "pedro-go/domain"
	"pedro-go/domain/expect"
	"pedro-go/ra"
	"testing"
)

func NewInMemoryArtistRegistry(raArtists map[ra.Slug]ra.ArtistWithEvents) *ArtistRegistry {
	repo := db.NewInMemoryArtistRepository()
	raClient := ra.NewInMemoryClient(raArtists)

	return NewArtistRegistry(repo, raClient)
}

func TestArtistRegistry(t *testing.T) {
	t.Run("lists all artists", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{},
				},
			},
		)

		err := registry.Follow("boysnoize", UserID(1))
		err = registry.Follow("sinamin", UserID(1))

		got := registry.All()
		want := Artists{
			{
				ID:         1,
				RAID:       "943",
				RASlug:     "boysnoize",
				Name:       "Boys Noize",
				FollowedBy: UserIDs{UserID(1)},
			},
			{
				ID:         2,
				RAID:       "222",
				RASlug:     "sinamin",
				Name:       "Sinamin",
				FollowedBy: UserIDs{UserID(1)},
			},
		}

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("follows an artist from resident advisor", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"daftpunk": {
					Artist:     ra.Artist{RAID: "111", Name: "Daft Punk"},
					EventsData: []ra.Event{},
				},
			},
		)

		want := []Artist{{ID: 1, RAID: "111", RASlug: "daftpunk", Name: "Daft Punk", FollowedBy: UserIDs{UserID(1)}}}
		err := registry.Follow("daftpunk", UserID(1))

		expect.NoErr(t, err)
		expect.DeepEqual(t, registry.All(), want)
	})

	t.Run("doesn't add artist if already added", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
			},
		)

		want := Artists{
			{ID: 1, RAID: "943", RASlug: "boysnoize", Name: "Boys Noize", FollowedBy: UserIDs{UserID(1)}},
		}
		err := registry.Follow("boysnoize", UserID(1))

		expect.NoErr(t, err)
		expect.DeepEqual(t, registry.All(), want)
	})

	t.Run("returns error if artist can't be found on RA", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{},
		)

		err := registry.Follow("unknown", UserID(1))

		expect.Err(t, err)
		expect.Equal(t, err.Error(), ErrNotFoundOnRA.Error())
	})

	t.Run("follows new artist as user", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{},
				},
			},
		)

		err := registry.Follow("boysnoize", UserID(1))
		err = registry.Follow("sinamin", UserID(2))

		got, err := registry.ArtistsFor(UserID(1))
		want := Artists{
			Artist{
				ID:         1,
				RAID:       "943",
				RASlug:     "boysnoize",
				Name:       "Boys Noize",
				FollowedBy: UserIDs{UserID(1)},
			},
		}

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("ignores follow if already following", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{},
				},
			},
		)
		err := registry.Follow("boysnoize", UserID(1))
		err = registry.Follow("boysnoize", UserID(1))

		got, err := registry.ArtistsFor(UserID(1))

		expect.NoErr(t, err)
		expect.True(t, len(got) == 1)
		expect.DeepEqual(t, got[0].FollowedBy, UserIDs{UserID(1)})
	})

	t.Run("follows existing artist", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{},
				},
			},
		)
		err := registry.Follow("boysnoize", UserID(1))
		err = registry.Follow("boysnoize", UserID(2))

		got, err := registry.ArtistsFor(UserID(2))

		expect.NoErr(t, err)
		expect.Equal(t, len(got), 1)
		//expect.DeepEqual(t, got[0].FollowedBy, UserIDs{UserID(1), UserID(2)})
	})

	t.Run("fetches all events for artist in the next month", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
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
							Date:       "2023-11-04T00:00:00.000",
							StartTime:  "2023-11-04T13:00:00.000",
							ContentUrl: "/events/1789025",
						},
						{
							Id:         "2",
							Title:      "Klubnacht 2",
							Date:       "2023-11-04T00:00:00.000",
							StartTime:  "2023-11-04T13:00:00.000",
							ContentUrl: "/events/1789025",
						},
					},
				},
			},
		)

		events, err := registry.AllEventsForArtist(Artist{
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
		events := []ra.Event{
			{
				Id:         "3",
				Title:      "Kater Blau Night",
				Date:       "2023-11-04T00:00:00.000",
				StartTime:  "2023-11-04T13:00:00.000",
				ContentUrl: "/events/3",
			},
			{
				Id:         "1",
				Title:      "Klubnacht",
				Date:       "2023-11-04T00:00:00.000",
				StartTime:  "2023-11-04T13:00:00.000",
				ContentUrl: "/events/1789025",
			},
			{
				Id:         "2",
				Title:      "Klubnacht 2",
				Date:       "2023-11-04T00:00:00.000",
				StartTime:  "2023-11-04T13:00:00.000",
				ContentUrl: "/events/1789025",
			},
		}
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Event{events[0]},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Event{events[1], events[2]},
				},
			},
		)

		joe := UserID(1)

		var err error
		err = registry.Follow("boysnoize", joe)
		err = registry.Follow("sinamin", joe)
		eventsForUser, err := registry.NewEventsForUser(joe)

		expect.NoErr(t, err)
		expect.DeepEqual(t, eventsForUser, events)

		newlyFetched, err := registry.NewEventsForUser(joe)

		expect.NoErr(t, err)
		expect.Len(t, newlyFetched, 0)
	})
}
