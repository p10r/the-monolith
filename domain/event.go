package domain

import (
	"slices"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Id         EventID
	Title      string
	Artist     string
	Venue      string
	City       string
	StartTime  time.Time
	ContentUrl string
}

type Events []Event

type EventID int64

func (events Events) IDs() EventIDs {
	ids := EventIDs{}
	for _, e := range events {
		ids = append(ids, e.Id)
	}
	return ids
}

func (events Events) FindNewEvents(a Artist) Events {
	// TODO there are artists being created without empty list - maybe in DB?
	if a.TrackedEvents == nil {
		a.TrackedEvents = EventIDs{}
	}

	newEvents := Events{}
	for _, e := range events {
		if strings.ToLower(e.City) != "berlin" {
			continue
		}

		if a.TrackedEvents.Contains(e.Id) {
			continue
		}
		newEvents = append(newEvents, e)
	}
	return newEvents
}

type EventIDs []EventID

func NewEventID(id string) (EventID, error) {
	i, err := strconv.ParseInt(strings.TrimSpace(id), 10, 64)
	if err != nil {
		return EventID(0), err
	}
	return EventID(i), nil
}

func (ids EventIDs) Contains(id EventID) bool {
	var ints []int64
	for _, eventID := range ids {
		ints = append(ints, int64(eventID))
	}

	return slices.Contains(ints, int64(id))
}
