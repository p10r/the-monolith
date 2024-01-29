package domain

import (
	"slices"
	"strconv"
	"strings"
)

type Event struct {
	Id         string
	Title      string
	Venue      string
	Date       string
	StartTime  string
	ContentUrl string
}

type Events []Event

type EventID int64

type EventIDs []EventID

// NewEventID should be replaced with this:
// https://stackoverflow.com/questions/51923863/how-to-construct-json-so-i-can-receive-int64-and-string-using-golang
func NewEventID(id string) (EventID, error) {
	i, err := strconv.ParseInt(strings.TrimSpace(id), 10, 64)
	if err != nil {
		return EventID(0), err
	}
	return EventID(i), nil
}

func (eventId EventIDs) Contains(id EventID) bool {
	var ints []int64
	for _, eventID := range eventId {
		ints = append(ints, int64(eventID))
	}

	return slices.Contains(ints, int64(id))
}
