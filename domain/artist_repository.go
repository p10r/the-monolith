package domain

import (
	"pedro-go/domain/expect"
	"testing"
)

type ArtistRepository interface {
	Save(artist Artist) (Artist, error)
	All() (Artists, error)
}

type ArtistRepositoryContract struct {
	NewRepository func() ArtistRepository
}

func (c ArtistRepositoryContract) Test(t *testing.T) {
	t.Run("adds artist", func(t *testing.T) {
		r := c.NewRepository()
		artist := Artist{
			RAId:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    UserIds{UserId(1)},
			TrackedEvents: EventIds{EventId(1)},
		}
		_, err := r.Save(artist)

		want := Artists{
			Artist{
				ID:            1,
				RAId:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    UserIds{UserId(1)},
				TrackedEvents: EventIds{EventId(1)},
			},
		}
		got, err := r.All()

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("updates artist", func(t *testing.T) {
		r := c.NewRepository()

		artist := Artist{
			RAId:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    UserIds{UserId(1)},
			TrackedEvents: EventIds{},
		}
		first, err := r.Save(artist)

		first.FollowedBy = UserIds{UserId(1), UserId(2)}
		_, err = r.Save(first)

		want := Artists{
			Artist{
				ID:            1,
				RAId:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    UserIds{UserId(1), UserId(2)},
				TrackedEvents: EventIds{},
			},
		}

		got, err := r.All()

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("updates events", func(t *testing.T) {
		r := c.NewRepository()

		artist := Artist{
			RAId:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    UserIds{},
			TrackedEvents: EventIds{EventId(1)},
		}
		first, err := r.Save(artist)

		first.TrackedEvents = EventIds{EventId(1), EventId(2)}
		_, err = r.Save(first)

		want := Artists{
			Artist{
				ID:            1,
				RAId:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    UserIds{},
				TrackedEvents: EventIds{EventId(1), EventId(2)},
			},
		}

		got, err := r.All()

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})
}
