package ra

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Event struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	StartTime  string `json:"startTime"`
	ContentUrl string `json:"contentUrl"`
	Images     []struct {
		Filename string `json:"filename"`
	} `json:"images"`
	Venue struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		ContentUrl string `json:"contentUrl"`
		Area       struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"area"`
	} `json:"venue"`
}

// TODO write tests
func NewEvent(res *http.Response, err error) ([]Event, error) {
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
