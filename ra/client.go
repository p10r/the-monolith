package ra

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

// TODO add pointer
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

func (c Client) GetEventsByArtistId(raId string, start time.Time, end time.Time) ([]Events, error) {
	req, err := GetEventsFor(raId, start, end, c.baseUri)

	res, err := c.http.Do(req)
	if err != nil {
		return []Events{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		var errRes ErrorRes
		if err = json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return []Events{}, err
		}

		return []Events{}, fmt.Errorf("Request failed with 400 - %v\n\n", errRes) //TODO formatting
	}

	if res.StatusCode != http.StatusOK {
		return []Events{}, fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	var events struct {
		Data struct {
			Listing struct {
				EventsData   []Events `json:"data"`
				TotalResults int      `json:"totalResults"`
			} `json:"listing"`
		} `json:"data"`
	}

	if err = json.NewDecoder(res.Body).Decode(&events); err != nil {
		log.Println("Can not deserialize response to EventsResponse")
		return []Events{}, err
	}
	return events.Data.Listing.EventsData, err
}
