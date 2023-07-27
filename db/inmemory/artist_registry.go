package inmemory

import "pedro-go/domain"

type ArtistRegistry struct {
	Artists domain.Artists
}

func (r ArtistRegistry) FindAll() domain.Artists {
	return domain.Artists{domain.Artist{Name: "Boys Noize"}}
}
