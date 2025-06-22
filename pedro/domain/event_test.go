package domain

import (
	"fmt"
	"github.com/p10r/monolith/pedro/domain/expect"
	"testing"
)

func TestEvents(t *testing.T) {

	testCases := []struct {
		Artist    Artist
		NewEvents EventIDs
		Want      EventIDs
	}{
		{
			Artist:    Artist{TrackedEvents: EventIDs{}},
			NewEvents: EventIDs{1, 2},
			Want:      EventIDs{1, 2},
		},
		{
			Artist:    Artist{TrackedEvents: EventIDs{1, 2}},
			NewEvents: EventIDs{1, 2, 3},
			Want:      EventIDs{3},
		},
		{
			Artist:    Artist{TrackedEvents: EventIDs{1, 2}},
			NewEvents: EventIDs{1, 2},
			Want:      EventIDs{},
		},
		{
			Artist:    Artist{TrackedEvents: EventIDs{1, 2}},
			NewEvents: EventIDs{},
			Want:      EventIDs{},
		},
		{
			Artist:    Artist{TrackedEvents: EventIDs{1, 2}},
			NewEvents: EventIDs{1, 2},
			Want:      EventIDs{},
		},
	}
	t.Run("filters out already tracked events", func(t *testing.T) {
		for _, tc := range testCases {
			got := tc.NewEvents.ToEvents().FindNewEvents(tc.Artist)
			expect.DeepEqual(t, got, tc.Want.ToEvents())
		}
	})

	//TODO fuzz test
	for id, input := range []struct {
		Artist  Artist
		City    string
		Matches bool
	}{
		{
			Artist:  Artist{},
			City:    "BERLIN",
			Matches: true,
		},
		{
			Artist:  Artist{},
			City:    "berlin",
			Matches: true,
		},
		{
			Artist:  Artist{},
			City:    "berlin",
			Matches: true,
		},
		{
			Artist:  Artist{},
			City:    "Ham burg",
			Matches: false,
		},
		{
			Artist:  Artist{},
			City:    "HAMBURG",
			Matches: false,
		},
	} {
		testcase := fmt.Sprintf("filters out cities except Berlin: %v", input.City)

		events := Events{Event{Id: EventID(id), City: input.City}}
		t.Run(testcase, func(t *testing.T) {
			if input.Matches {
				expect.Len(t, events.FindEventsInBerlin(Artist{}), 1)
			} else {
				expect.Len(t, events.FindEventsInBerlin(Artist{}), 0)
			}
		})
	}
}

func (ids EventIDs) ToEvents() Events {
	es := Events{}
	for _, id := range ids {
		es = append(es, Event{Id: id, City: "Berlin"})
	}
	return es
}
