package domain

import (
	"pedro-go/pedro/domain/expect"
	"testing"
)

func TestArtist(t *testing.T) {
	t.Run("add follower", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAID:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIDs{},
		}

		got := artist.AddFollower(UserID(1))

		want := artist
		want.FollowedBy = UserIDs{UserID(1)}

		expect.DeepEqual(t, got, want)
	})

	t.Run("ignores if already following", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAID:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIDs{},
		}

		got := artist.AddFollower(UserID(1)).AddFollower(UserID(1))

		want := artist
		want.FollowedBy = UserIDs{UserID(1)}

		expect.DeepEqual(t, got, want)
	})

	t.Run("adds another follower on top", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAID:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIDs{},
		}

		got := artist.AddFollower(UserID(1)).AddFollower(UserID(2))

		want := artist
		want.FollowedBy = UserIDs{UserID(1), UserID(2)}

		expect.DeepEqual(t, got, want)
	})

	t.Run("removes follower", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAID:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIDs{UserID(1)},
		}

		got := artist.RemoveFollower(UserID(1))
		want := artist
		want.FollowedBy = UserIDs{}

		expect.DeepEqual(t, got, want)
	})

	t.Run("doesn't remove follower if not present", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAID:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIDs{},
		}

		got := artist.RemoveFollower(UserID(1))

		want := artist
		expect.DeepEqual(t, got, want)
	})

	t.Run("filters artists by user id", func(t *testing.T) {
		artists := Artists{
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
				FollowedBy: UserIDs{UserID(1), UserID(2)},
			},
			{
				ID:         3,
				RAID:       "111",
				RASlug:     "daftpunk",
				Name:       "Daft Punk",
				FollowedBy: UserIDs{UserID(3)},
			},
		}

		got := artists.FilterByUserId(UserID(1))
		want := Artists{artists[0], artists[1]}

		expect.DeepEqual(t, got, want)
	})

	t.Run("returns empty list if none is matching the user id", func(t *testing.T) {
		artists := Artists{
			{
				ID:         1,
				RAID:       "943",
				RASlug:     "boysnoize",
				Name:       "Boys Noize",
				FollowedBy: UserIDs{UserID(1)},
			},
		}

		got := len(artists.FilterByUserId(UserID(2)))

		expect.Equal(t, got, 0)
	})
}
