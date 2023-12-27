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
					EventsData: []ra.Events{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Events{},
				},
			},
		)

		err := registry.Follow("boysnoize", UserId(1))
		err = registry.Follow("sinamin", UserId(1))

		got := registry.All()
		want := Artists{
			{
				Id:         1,
				RAId:       "943",
				RASlug:     "boysnoize",
				Name:       "Boys Noize",
				FollowedBy: UserIds{UserId(1)},
			},
			{
				Id:         2,
				RAId:       "222",
				RASlug:     "sinamin",
				Name:       "Sinamin",
				FollowedBy: UserIds{UserId(1)},
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
					EventsData: []ra.Events{},
				},
			},
		)

		want := []Artist{{Id: 1, RAId: "111", RASlug: "daftpunk", Name: "Daft Punk", FollowedBy: UserIds{UserId(1)}}}
		err := registry.Follow("daftpunk", UserId(1))

		expect.NoErr(t, err)
		expect.DeepEqual(t, registry.All(), want)
	})

	t.Run("doesn't add artist if already added", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Events{},
				},
			},
		)

		want := Artists{
			{Id: 1, RAId: "943", RASlug: "boysnoize", Name: "Boys Noize", FollowedBy: UserIds{UserId(1)}},
		}
		err := registry.Follow("boysnoize", UserId(1))

		expect.NoErr(t, err)
		expect.DeepEqual(t, registry.All(), want)
	})

	t.Run("returns error if artist can't be found on RA", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{},
		)

		err := registry.Follow("unknown", UserId(1))

		expect.Err(t, err)
		expect.Equal(t, err.Error(), ErrNotFoundOnRA.Error())
	})

	//t.Run("adds slug to queue if RA is not reachable", func(t *testing.T) {
	//	t.Fail() TODO
	//})

	t.Run("follows new artist as user", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Events{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Events{},
				},
			},
		)

		err := registry.Follow("boysnoize", UserId(1))
		err = registry.Follow("sinamin", UserId(2))

		got, err := registry.ArtistsFor(UserId(1))
		want := Artists{
			Artist{
				Id:         1,
				RAId:       "943",
				RASlug:     "boysnoize",
				Name:       "Boys Noize",
				FollowedBy: UserIds{UserId(1)},
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
					EventsData: []ra.Events{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Events{},
				},
			},
		)
		err := registry.Follow("boysnoize", UserId(1))
		err = registry.Follow("boysnoize", UserId(1))

		got, err := registry.ArtistsFor(UserId(1))

		expect.NoErr(t, err)
		expect.True(t, len(got) == 1)
		expect.DeepEqual(t, got[0].FollowedBy, UserIds{UserId(1)})
	})

	t.Run("follows existing artist", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Events{},
				},
				"sinamin": {
					Artist:     ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Events{},
				},
			},
		)
		err := registry.Follow("boysnoize", UserId(1))
		err = registry.Follow("boysnoize", UserId(2))

		got, err := registry.ArtistsFor(UserId(2))

		expect.NoErr(t, err)
		expect.Equal(t, len(got), 1)
		//expect.DeepEqual(t, got[0].FollowedBy, UserIds{UserId(1), UserId(2)})
	})

	t.Run("fetches all events for artist in the next 7 days", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.ArtistWithEvents{
				"boysnoize": {
					Artist:     ra.Artist{RAID: "943", Name: "Boys Noize"},
					EventsData: []ra.Events{},
				},
				"sinamin": {
					Artist: ra.Artist{RAID: "222", Name: "Sinamin"},
					EventsData: []ra.Events{
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

		events, err := registry.EventsFor(Artist{
			Id:         1,
			RAId:       "222",
			RASlug:     "sinamin",
			Name:       "Sinamin",
			FollowedBy: nil,
		})

		expect.NoErr(t, err)
		expect.Equal(t, len(events), 2)
	})

}
