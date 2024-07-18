package domain

import (
	"context"
)

type ArtistRepository interface {
	Save(ctx context.Context, artist Artist) (Artist, error)
	All(ctx context.Context) (Artists, error)
}
