package db

import (
	"fmt"
	"pedro-go/domain"
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

func (r *InMemoryArtistRepository) Add(artist domain.Artist) {
	artist.Id = r.id.Load()
	r.id.Add(1)

	fmt.Printf("InMemoryArtistRepository: Adding %v\n", artist)
	r.artists[artist.Id] = artist
}

func (r *InMemoryArtistRepository) All() []domain.Artist {
	fmt.Println("InMemoryArtistRepository: Requesting all artists")

	a := make([]domain.Artist, 0, len(r.artists))
	for _, artist := range r.artists {
		a = append(a, artist)
	}
	return a
}
