package domain_test

import (
	"pedro-go/db/inmemory"
	. "pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestArtistRegistry(t *testing.T) {
	t.Run("lists all artists", func(t *testing.T) {
		repo := inmemory.NewInMemoryArtistRepository()
		registry := NewArtistRegistry(repo)

		repo.Add(Artist{RAId: 943, Name: "A"})
		repo.Add(Artist{RAId: 222, Name: "B"})

		expect.SliceContains(
			t, registry.All(),
			Artist{Id: 1, RAId: 943, Name: "A"},
			Artist{Id: 2, RAId: 222, Name: "B"},
		)
	})
}
