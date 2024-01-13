package domain

import (
	"pedro-go/domain/expect"
	"testing"
)

func TestArtist(t *testing.T) {
	t.Run("add follower", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAId:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIds{},
		}

		got := artist.AddFollower(UserId(1))

		want := artist
		want.FollowedBy = UserIds{UserId(1)}

		expect.DeepEqual(t, got, want)
	})

	t.Run("ignores if already following", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAId:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIds{},
		}

		got := artist.AddFollower(UserId(1)).AddFollower(UserId(1))

		want := artist
		want.FollowedBy = UserIds{UserId(1)}

		expect.DeepEqual(t, got, want)
	})

	t.Run("adds another follower on top", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAId:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIds{},
		}

		got := artist.AddFollower(UserId(1)).AddFollower(UserId(2))

		want := artist
		want.FollowedBy = UserIds{UserId(1), UserId(2)}

		expect.DeepEqual(t, got, want)
	})

	t.Run("removes follower", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAId:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIds{UserId(1)},
		}

		got := artist.RemoveFollower(UserId(1))
		want := artist
		want.FollowedBy = UserIds{}

		expect.DeepEqual(t, got, want)
	})

	t.Run("doesn't remove follower if not present", func(t *testing.T) {
		artist := Artist{
			ID:         1,
			RAId:       "943",
			RASlug:     "boysnoize",
			Name:       "Boys Noize",
			FollowedBy: UserIds{},
		}

		got := artist.RemoveFollower(UserId(1))

		want := artist
		expect.DeepEqual(t, got, want)
	})

	t.Run("filters artists by user id", func(t *testing.T) {
		artists := Artists{
			{
				ID:         1,
				RAId:       "943",
				RASlug:     "boysnoize",
				Name:       "Boys Noize",
				FollowedBy: UserIds{UserId(1)},
			},
			{
				ID:         2,
				RAId:       "222",
				RASlug:     "sinamin",
				Name:       "Sinamin",
				FollowedBy: UserIds{UserId(1), UserId(2)},
			},
			{
				ID:         3,
				RAId:       "111",
				RASlug:     "daftpunk",
				Name:       "Daft Punk",
				FollowedBy: UserIds{UserId(3)},
			},
		}

		got := artists.FilterByUserId(UserId(1))
		want := Artists{artists[0], artists[1]}

		expect.DeepEqual(t, got, want)
	})

	t.Run("returns empty list if none is matching the user id", func(t *testing.T) {
		artists := Artists{
			{
				ID:         1,
				RAId:       "943",
				RASlug:     "boysnoize",
				Name:       "Boys Noize",
				FollowedBy: UserIds{UserId(1)},
			},
		}

		got := len(artists.FilterByUserId(UserId(2)))

		expect.Equal(t, got, 0)
	})
}
