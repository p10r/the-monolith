package domain

import (
	"slices"
)

type Artist struct {
	ID            int64
	RAID          string
	RASlug        RASlug
	Name          string
	FollowedBy    UserIDs
	TrackedEvents EventIDs
}

// ArtistInfo is the ra.co representation of an artist when searching by slug
type ArtistInfo struct {
	RAID string
	Name string
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

func (a Artists) RASlugs() []RASlug {
	slugs := make([]RASlug, len(a))
	for _, artist := range a {
		slugs = append(slugs, artist.RASlug)
	}
	return slugs
}

type UserID int64

type UserIDs []UserID
