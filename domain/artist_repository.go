package domain

import (
	"pedro-go/domain/expect"
	"testing"
)

type ArtistRepository interface {
	Save(artist Artist) Artist
	All() []Artist
}

type ArtistRepositoryContract struct {
	NewRepository func() ArtistRepository
}

func (c ArtistRepositoryContract) Test(t *testing.T) {
	t.Run("returns all artists", func(t *testing.T) {

	})

	t.Run("adds artist", func(t *testing.T) {
		r := c.NewRepository()
		r.Save(Artist{RAId: "943", RASlug: "boysnoize", Name: "Boys Noize"})

		want := Artists{Artist{Id: 1, RAId: "943", RASlug: "boysnoize", Name: "Boys Noize"}}
		expect.DeepEqual(t, r.All(), want)
	})

	t.Run("updates artist", func(t *testing.T) {
		r := c.NewRepository()

		stored := r.Save(Artist{RAId: "943", RASlug: "boysnoize", Name: "Boys Noize", FollowedBy: UserIds{UserId(1)}})
		stored.FollowedBy = UserIds{UserId(1), UserId(2)}
		r.Save(stored)

		want := Artists{
			Artist{
				Id:         1,
				RAId:       "943",
				RASlug:     "boysnoize",
				Name:       "Boys Noize",
				FollowedBy: UserIds{UserId(1), UserId(2)},
			},
		}
		expect.DeepEqual(t, r.All(), want)
	})

}
