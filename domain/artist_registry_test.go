package domain_test

import (
	db "pedro-go/db/inmemory"
	. "pedro-go/domain"
	"pedro-go/domain/expect"
	"pedro-go/ra"
	"pedro-go/ra/inmemory"
	"testing"
)

func NewInMemoryArtistRegistry(raArtists map[ra.Slug]ra.Artist) *ArtistRegistry {
	repo := db.NewInMemoryArtistRepository()
	raClient := inmemory.NewClient(raArtists)

	return NewArtistRegistry(repo, raClient)
}

func TestArtistRegistry(t *testing.T) {
	t.Run("lists all artists", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.Artist{
				"boysnoize": {RAID: "943", Name: "Boys Noize"},
				"sinamin":   {RAID: "222", Name: "Sinamin"},
			},
		)

		err := registry.Add("boysnoize")
		err = registry.Add("sinamin")

		expect.NoErr(t, err)
		expect.SliceContains(
			t, registry.All(),
			Artist{Id: 1, RAId: "943", RASlug: "boysnoize", Name: "Boys Noize"},
			Artist{Id: 2, RAId: "222", RASlug: "sinamin", Name: "Sinamin"},
		)
	})

	t.Run("adds an artist from resident advisor", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.Artist{
				"daftpunk": {RAID: "111", Name: "Daft Punk"},
			},
		)

		want := Artist{Id: 1, RAId: "111", RASlug: "daftpunk", Name: "Daft Punk"}
		err := registry.Add("daftpunk")

		expect.NoErr(t, err)
		expect.SliceContains(t, registry.All(), want)
	})

	t.Run("doesn't add artist if already added", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.Artist{
				"boysnoize": {RAID: "943", Name: "Boys Noize"},
			},
		)

		want := Artists{Artist{Id: 1, RAId: "943", RASlug: "boysnoize", Name: "Boys Noize"}}
		err := registry.Add("boysnoize")

		expect.NoErr(t, err)
		expect.SliceEqual(t, registry.All(), want)
	})

	t.Run("returns error if artist can't be found on RA", func(t *testing.T) {
		registry := NewInMemoryArtistRegistry(
			map[ra.Slug]ra.Artist{},
		)
		err := registry.Add("unknown")

		expect.Err(t, err)
		expect.Equal(t, err.Error(), ErrNotFoundOnRA.Error())
	})

	//t.Run("adds slug to queue if RA is not reachable", func(t *testing.T) {
	//	t.Fail() TODO
	//})
}
