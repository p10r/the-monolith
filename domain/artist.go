package domain

import (
	"pedro-go/ra"
	"slices"
	"strconv"
	"strings"
)

type Artist struct {
	ID            int64
	RAID          string
	RASlug        ra.Slug
	Name          string
	FollowedBy    UserIDs
	TrackedEvents EventIDs
}

func (a Artist) AddFollower(id UserID) Artist {
	if slices.Contains(a.FollowedBy, id) {
		return a
	}
	a.FollowedBy = append(a.FollowedBy, id)
	return a
}

func (a Artist) RemoveFollower(id UserID) Artist {
	a.FollowedBy = slices.DeleteFunc(a.FollowedBy, func(existingId UserID) bool {
		return existingId == id
	})
	return a
}

type Artists []Artist

func (a Artists) FilterByUserId(id UserID) Artists {
	var found Artists
	for _, artist := range a {
		if slices.Contains(artist.FollowedBy, id) {
			found = append(found, artist)
		}
	}
	return found
}

func (a Artists) RASlugs() []ra.Slug {
	slugs := make([]ra.Slug, len(a))
	for _, artist := range a {
		slugs = append(slugs, artist.RASlug)
	}
	return slugs
}

type UserID int64
type UserIDs []UserID

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
