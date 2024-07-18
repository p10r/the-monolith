package db

import (
	"context"
	d "github.com/p10r/pedro/pedro/domain"
	"github.com/p10r/pedro/pedro/domain/expect"
	"github.com/p10r/pedro/pkg/sqlite"
	"testing"
)

func NewRepository(t *testing.T) *SqliteArtistRepository {
	conn := sqlite.MustOpenDB(t)
	return NewSqliteArtistRepository(conn)
}

func TestSqliteArtistRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("adds artist", func(t *testing.T) {
		r := NewRepository(t)
		artist := d.Artist{
			RAID:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    d.UserIDs{d.UserID(1)},
			TrackedEvents: d.EventIDs{d.EventID(1)},
		}
		_, err := r.Save(ctx, artist)
		expect.NoErr(t, err)

		want := d.Artists{
			d.Artist{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    d.UserIDs{d.UserID(1)},
				TrackedEvents: d.EventIDs{d.EventID(1)},
			},
		}
		got, err := r.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("updates artist", func(t *testing.T) {
		r := NewRepository(t)

		artist := d.Artist{
			RAID:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    d.UserIDs{d.UserID(1)},
			TrackedEvents: d.EventIDs{},
		}
		first, err := r.Save(ctx, artist)
		expect.NoErr(t, err)

		first.FollowedBy = d.UserIDs{d.UserID(1), d.UserID(2)}
		_, err = r.Save(ctx, first)
		expect.NoErr(t, err)

		want := d.Artists{
			d.Artist{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    d.UserIDs{d.UserID(1), d.UserID(2)},
				TrackedEvents: d.EventIDs{},
			},
		}

		got, err := r.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	t.Run("updates events", func(t *testing.T) {
		r := NewRepository(t)

		artist := d.Artist{
			RAID:          "943",
			RASlug:        "boysnoize",
			Name:          "Boys Noize",
			FollowedBy:    d.UserIDs{},
			TrackedEvents: d.EventIDs{d.EventID(1)},
		}
		first, err := r.Save(ctx, artist)
		expect.NoErr(t, err)

		first.TrackedEvents = d.EventIDs{d.EventID(1), d.EventID(2)}
		_, err = r.Save(ctx, first)
		expect.NoErr(t, err)

		want := d.Artists{
			d.Artist{
				ID:            1,
				RAID:          "943",
				RASlug:        "boysnoize",
				Name:          "Boys Noize",
				FollowedBy:    d.UserIDs{},
				TrackedEvents: d.EventIDs{d.EventID(1), d.EventID(2)},
			},
		}

		got, err := r.All(ctx)

		expect.NoErr(t, err)
		expect.DeepEqual(t, got, want)
	})

	//same is being mapped for domain.EventIDs
	t.Run("map domain IDs to string list", func(t *testing.T) {
		for _, tc := range []struct {
			Input d.UserIDs
			Want  commaSeparatedStr
		}{
			{
				Input: d.UserIDs{d.UserID(1), d.UserID(2), d.UserID(3), d.UserID(4)},
				Want:  commaSeparatedStr("1,2,3,4"),
			},
			{
				Input: d.UserIDs{},
				Want:  commaSeparatedStr(""),
			},
		} {
			got := newUserIdString(tc.Input)
			expect.Equal(t, got, tc.Want)
		}
	})

	//same is being mapped for domain.EventIDs
	t.Run("map id string to domain ids", func(t *testing.T) {
		for _, tc := range []struct {
			Input commaSeparatedStr
			Want  d.UserIDs
		}{
			{
				Input: commaSeparatedStr("1,2,3,4"),
				Want:  d.UserIDs{d.UserID(1), d.UserID(2), d.UserID(3), d.UserID(4)},
			},
			{
				Input: commaSeparatedStr(""),
				Want:  d.UserIDs{},
			},
			{
				Input: commaSeparatedStr("1 , 2 , 3 ,     4     "),
				Want:  d.UserIDs{d.UserID(1), d.UserID(2), d.UserID(3), d.UserID(4)},
			},
		} {
			expect.DeepEqual(t, tc.Input.toUserIds(), tc.Want)
		}
	})
}
