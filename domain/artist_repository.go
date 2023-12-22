package domain

type ArtistRepository interface {
	Add(artist Artist)
	All() []Artist
}
