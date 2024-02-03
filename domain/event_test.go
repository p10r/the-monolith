package domain

import (
	"pedro-go/domain/expect"
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
}

func (ids EventIDs) ToEvents() Events {
	es := Events{}
	for _, id := range ids {
		es = append(es, Event{Id: id})
	}
	return es
}
