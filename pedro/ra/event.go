package ra

import (
	"encoding/json"
	"fmt"
	"github.com/p10r/pedro/pedro/domain"
	"log"
	"net/http"
	"time"
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
	Venue Venue `json:"venue"`
}

type Venue struct {
	Name string `json:"name"`
	Area Area   `json:"area"`
}

type Area struct {
	Name string `json:"name"`
}

type Events []Event

func NewEvent(res *http.Response, err error) (Events, error) {
	if res == nil {
		return Events{}, fmt.Errorf("ra events response is nil")
	}

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

func (events Events) ToDomainEvents(artistName string) domain.Events {
	d := domain.Events{}
	for _, e := range events {
		id, err := domain.NewEventID(e.Id)
		if err != nil {
			log.Printf("failed parsing %v to EventID: %v", e.Id, err)
			continue
		}

		layout := "2006-01-02T15:04:05.000"
		date, err := time.Parse(layout, e.StartTime)

		if err != nil {
			log.Printf("failed parsing %v to time: %v", e.Date, err)
			continue
		}

		transformed := domain.Event{
			Id:         id,
			Title:      e.Title,
			Artist:     artistName,
			Venue:      e.Venue.Name,
			City:       e.Venue.Area.Name,
			StartTime:  date,
			ContentUrl: e.ContentUrl,
		}

		d = append(d, transformed)
	}
	return d
}
