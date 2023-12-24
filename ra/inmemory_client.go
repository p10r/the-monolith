package ra

import (
	"errors"
)

type InMemoryClient struct {
	artists map[Slug]Artist
}

func NewInMemoryClient(artists map[Slug]Artist) *InMemoryClient {
	return &InMemoryClient{artists: artists}
}

func (c *InMemoryClient) GetArtistBySlug(slug Slug) (Artist, error) {
	res, ok := c.artists[slug]

	if ok {
		return res, nil
	}

	return Artist{}, errors.New("artist not found")
}
