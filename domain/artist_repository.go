package domain

import (
	"context"
	"pedro-go/domain/expect"
	"testing"
)

type ArtistRepository interface {
	Save(ctx context.Context, artist Artist) (Artist, error)
	All(ctx context.Context) (Artists, error)
}

type ArtistRepositoryContract struct {
	NewRepository func() ArtistRepository
}

func (c ArtistRepositoryContract) Test(t *testing.T) {
	ctx := context.Background()

	t.Run("adds artist", func(t *testing.T) {
		r := c.NewRepository()
		artist := Artist{
			RAID:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    UserIDs{UserID(1)},
			TrackedEvents: EventIDs{EventID(1)},
		}
		_, err := r.Save(ctx, artist)
		expect.NoErr(t, err)

		want := Artists{
			Artist{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    UserIDs{UserID(1)},
				TrackedEvents: EventIDs{EventID(1)},
			},
		}
		got, err := r.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("updates artist", func(t *testing.T) {
		r := c.NewRepository()

		artist := Artist{
			RAID:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    UserIDs{UserID(1)},
			TrackedEvents: EventIDs{},
		}
		first, err := r.Save(ctx, artist)
		expect.NoErr(t, err)

		first.FollowedBy = UserIDs{UserID(1), UserID(2)}
		_, err = r.Save(ctx, first)
		expect.NoErr(t, err)

		want := Artists{
			Artist{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    UserIDs{UserID(1), UserID(2)},
				TrackedEvents: EventIDs{},
			},
		}

		got, err := r.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("updates events", func(t *testing.T) {
		r := c.NewRepository()

		artist := Artist{
			RAID:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    UserIDs{},
			TrackedEvents: EventIDs{EventID(1)},
		}
		first, err := r.Save(ctx, artist)
		expect.NoErr(t, err)

		first.TrackedEvents = EventIDs{EventID(1), EventID(2)}
		_, err = r.Save(ctx, first)
		expect.NoErr(t, err)

		want := Artists{
			Artist{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    UserIDs{},
				TrackedEvents: EventIDs{EventID(1), EventID(2)},
			},
		}

		got, err := r.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})
}
