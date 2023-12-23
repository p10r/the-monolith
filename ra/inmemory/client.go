package inmemory

import (
	"pedro-go/ra"
)

type Client struct {
	artists map[ra.Slug]ra.Artist
}

func NewClient(artists map[ra.Slug]ra.Artist) *Client {
	return &Client{artists: artists}
}

func (c *Client) GetArtistBySlug(slug ra.Slug) (ra.Artist, error) {
	return c.artists[slug], nil
}
