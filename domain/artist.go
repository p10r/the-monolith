package domain

import (
	"pedro-go/ra"
	"slices"
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
