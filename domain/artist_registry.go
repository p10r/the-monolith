package domain

import "context"

type ArtistRegistry interface {
	Add(ctx context.Context, artist NewArtist) (Artist, error)
	FindAll(ctx context.Context) (Artists, error)
}
