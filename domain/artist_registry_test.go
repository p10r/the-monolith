package domain_test

import (
	db "pedro-go/db/inmemory"
	. "pedro-go/domain"
	"pedro-go/domain/expect"
	"pedro-go/ra"
	"pedro-go/ra/inmemory"
	"testing"
)

func TestArtistRegistry(t *testing.T) {
	repo := db.NewInMemoryArtistRepository()
	repo.Add(Artist{RAId: "943", RASlug: "boysnoize", Name: "Boys Noize"})
	repo.Add(Artist{RAId: "222", RASlug: "sinamin", Name: "Sinamin"})

	raArtists := map[ra.Slug]ra.Artist{
		ra.Slug("boysnoize"): {RAID: "943", Name: "Boys Noize"},
		ra.Slug("sinamin"):   {RAID: "222", Name: "Sinamin"},
		ra.Slug("daftpunk"):  {RAID: "111", Name: "Daft Punk"},
	}
	raClient := inmemory.NewClient(raArtists)

	registry := NewArtistRegistry(repo, raClient)

	t.Run("lists all artists", func(t *testing.T) {
		expect.SliceContains(
			t, registry.All(),
			Artist{Id: 1, RAId: "943", RASlug: "boysnoize", Name: "Boys Noize"},
			Artist{Id: 2, RAId: "222", RASlug: "sinamin", Name: "Sinamin"},
		)
	})

	t.Run("adds an artist from resident advisor", func(t *testing.T) {
		artist := Artist{Id: 3, RAId: "111", RASlug: "daftpunk", Name: "Daft Punk"}

		expect.SliceContainsNot(t, repo.All(), artist)

		err := registry.Add("daftpunk")

		expect.NoErr(t, err)
		expect.SliceContains(t, repo.All(), artist)
	})

	t.Run("doesn't add artist if already added", func(t *testing.T) {
		want := repo.All()

		err := registry.Add("boysnoize")

		expect.NoErr(t, err)
		expect.SliceEqual(t, repo.All(), want)
	})

	t.Run("returns error if artist can't be found on RA", func(t *testing.T) {
		err := registry.Add("unknown")

		expect.Err(t, err)
		expect.Equal(t, err.Error(), ErrNotFoundOnRA.Error())
	})

	//t.Run("adds slug to queue if RA is not reachable", func(t *testing.T) {
	//	t.Fail() TODO
	//})
}
