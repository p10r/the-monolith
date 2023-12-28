package ra

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

var ErrSlugNotFound = errors.New("slug not found on ra.co")

type Client struct {
	http    *http.Client
	baseUri string
}

func NewClient(baseUri string) *Client {
	return &Client{http: &http.Client{}, baseUri: baseUri}
}

func (c *Client) GetArtistBySlug(slug Slug) (Artist, error) {
	req, err := newGetArtistReq(slug, c.baseUri)
	if err != nil {
		return Artist{}, fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.http.Do(req)
	return NewArtist(res, err)
}

func (c *Client) GetEventsByArtistId(
	raId string,
	start time.Time,
	end time.Time,
) ([]Event, error) {
	req, err := newGetEvensReq(raId, start, end, c.baseUri)
	if err != nil {
		return []Event{}, fmt.Errorf("could not create request: %w", err)
	}

	res, err := c.http.Do(req)
	return NewEvent(res, err)
}
