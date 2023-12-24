package domain

import (
	"pedro-go/domain/expect"
	"pedro-go/ra"
	"testing"
)

type ResidentAdvisor interface {
	GetArtistBySlug(slug ra.Slug) (ra.Artist, error)
}

type RAContract struct {
	NewRA func() ResidentAdvisor
}

func (c RAContract) Test(t *testing.T) {
	client := c.NewRA()

	t.Run("returns artist by slug", func(t *testing.T) {
		artist, err := client.GetArtistBySlug("boysnoize")

		expect.NoErr(t, err)
		expect.DeepEqual(t, artist, ra.Artist{RAID: "943", Name: "Boys Noize"})
	})

	t.Run("returns ErrNotFound if slug can't be found", func(t *testing.T) {
		_, err := client.GetArtistBySlug("unknownabc")

		expect.Err(t, err)
		expect.DeepEqual(t, err, ra.ErrSlugNotFound)
	})
}
