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

func (c *Client) GetArtistBySlug(slug Slug) (Artist, error) {
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

func (c *Client) GetEventsByArtistId(raId string, start time.Time, end time.Time) ([]Event, error) {
	req, err := GetEventsFor(raId, start, end, c.baseUri)

	res, err := c.http.Do(req)
	if err != nil {
		return []Event{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusBadRequest {
		var errRes ErrorRes
		if err = json.NewDecoder(res.Body).Decode(&errRes); err != nil {
			return []Event{}, err
		}

		return []Event{}, fmt.Errorf("Request failed with 400 - %v\n\n", errRes) //TODO formatting
	}

	if res.StatusCode != http.StatusOK {
		return []Event{}, fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	var events struct {
		Data struct {
			Listing struct {
				EventsData   []Event `json:"data"`
				TotalResults int     `json:"totalResults"`
			} `json:"listing"`
		} `json:"data"`
	}

	if err = json.NewDecoder(res.Body).Decode(&events); err != nil {
		log.Println("Can not deserialize response to EventsResponse")
		return []Event{}, err
	}
	return events.Data.Listing.EventsData, err
}
