package inmemory

import (
	"errors"
	"pedro-go/ra"
)

type Client struct {
	artists map[ra.Slug]ra.Artist
}

func NewClient(artists map[ra.Slug]ra.Artist) *Client {
	return &Client{artists: artists}
}

func (c *Client) GetArtistBySlug(slug ra.Slug) (ra.Artist, error) {
	res, ok := c.artists[slug]

	if ok {
		return res, nil
	}

	return ra.Artist{}, errors.New("artist not found")
}
