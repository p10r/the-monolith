package db

import (
	"cmp"
	"context"
	"fmt"
	"github.com/p10r/pedro/pedro/domain"
	"slices"
	"sync/atomic"
)

type InMemoryArtistRepository struct {
	id      *atomic.Int64
	artists map[int64]domain.Artist
}

func NewInMemoryArtistRepository() *InMemoryArtistRepository {
	var id atomic.Int64
	id.Store(1)
	return &InMemoryArtistRepository{id: &id, artists: map[int64]domain.Artist{}}
}

func (r *InMemoryArtistRepository) Save(
	_ context.Context,
	artist domain.Artist,
) (domain.Artist, error) {
	if artist.ID != 0 {
		r.artists[artist.ID] = artist
		return artist, nil
	}

	artist.ID = r.id.Load()
	r.id.Add(1)

	r.artists[artist.ID] = artist

	return artist, nil
}

func (r *InMemoryArtistRepository) All(_ context.Context) (domain.Artists, error) {
	fmt.Println("InMemoryArtistRepository: Requesting all artists")

	a := make([]domain.Artist, 0, len(r.artists))
	for _, artist := range r.artists {
		a = append(a, artist)
	}

	//Order by ID
	slices.SortFunc(a, func(a, b domain.Artist) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return a, nil
}
