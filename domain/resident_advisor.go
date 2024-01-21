package domain

import (
	"errors"
	"pedro-go/domain/expect"
	"testing"
	"time"
)

type ResidentAdvisor interface {
	GetArtistBySlug(slug RASlug) (ArtistInfo, error)
	GetEventsByArtistId(raId string, start time.Time, end time.Time) (Events, error)
}

type RAContract struct {
	NewRA func() ResidentAdvisor
}

func (c RAContract) Test(t *testing.T) {
	client := c.NewRA()

	t.Run("returns artist by slug", func(t *testing.T) {
		artist, err := client.GetArtistBySlug("boysnoize")

		expect.NoErr(t, err)
		expect.DeepEqual(t, artist, ArtistInfo{RAID: "943", Name: "Boys Noize"})
	})

	t.Run("returns ErrNotFound if slug can't be found", func(t *testing.T) {
		_, err := client.GetArtistBySlug("unknownabc")

		expect.Err(t, err)
		expect.DeepEqual(t, err, errors.New("slug not found on ra.co"))
	})

	t.Run("returns events for artist", func(t *testing.T) {
		juneFirst23 := time.Date(2023, 11, 1, 0, 0, 0, 0, time.UTC)
		julyFirst23 := time.Date(2023, 11, 15, 0, 0, 0, 0, time.UTC)

		events, err := client.GetEventsByArtistId("106972", juneFirst23, julyFirst23)

		expect.NoErr(t, err)
		expect.Equal(t, len(events), 2)
	})
}
