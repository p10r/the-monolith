package db_test

import (
	"pedro-go/db"
	"pedro-go/domain"
	"testing"
)

func TestInMemoryArtistRepository(t *testing.T) {
	t.Run("verify contract for in-memory repo", func(t *testing.T) {
		domain.ArtistRepositoryContract{NewRepository: func() domain.ArtistRepository {
			return db.NewInMemoryArtistRepository()
		}}.Test(t)
	})
}
