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
	repo.Add(Artist{RAId: "943", Name: "A"})
	repo.Add(Artist{RAId: "222", Name: "B"})

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
			Artist{Id: 1, RAId: "943", Name: "A"},
			Artist{Id: 2, RAId: "222", Name: "B"},
		)
	})

	t.Run("adds an artist from resident advisor", func(t *testing.T) {
		expect.SliceContainsNot(t, repo.All(), Artist{Id: 3, RAId: "111", Name: "Daft Punk"})

		registry.Add("daftpunk")

		expect.SliceContains(t, repo.All(), Artist{Id: 3, RAId: "111", Name: "Daft Punk"})
	})

	t.Run("doesn't add artist if already added", func(t *testing.T) {
		t.Fail()
	})

	t.Run("returns error if artist can't be found on RA", func(t *testing.T) {
		t.Fail()
	})

	t.Run("adds slug to queue if RA is not reachable", func(t *testing.T) {
		t.Fail()
	})

}
