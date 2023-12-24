package ra

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrSlugNotFound = errors.New("slug not found on ra.co")

type Client struct {
	http    *http.Client
	baseUri string
}

func NewClient(baseUri string) *Client {
	return &Client{http: &http.Client{}, baseUri: baseUri}
}

func (c Client) GetArtistBySlug(slug Slug) (Artist, error) {
	req, err := getArtistBySlugReq(slug, c.baseUri)

	res, err := c.http.Do(req)
	if err != nil {
		return Artist{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Artist{}, fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	arist, err := NewArtistFrom(res.Body)
	if arist == (Artist{}) {
		return Artist{}, ErrSlugNotFound
	}

	return arist, err
}
