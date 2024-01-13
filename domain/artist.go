package domain

import (
	"pedro-go/ra"
	"slices"
)

type Artist struct {
	ID            int64
	RAId          string
	RASlug        ra.Slug
	Name          string
	FollowedBy    UserIds
	TrackedEvents EventIds
}

func (a Artist) AddFollower(id UserId) Artist {
	if slices.Contains(a.FollowedBy, id) {
		return a
	}
	a.FollowedBy = append(a.FollowedBy, id)
	return a
}

func (a Artist) RemoveFollower(id UserId) Artist {
	a.FollowedBy = slices.DeleteFunc(a.FollowedBy, func(existingId UserId) bool {
		return existingId == id
	})
	return a
}

type Artists []Artist

func (a Artists) FilterByUserId(id UserId) Artists {
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

type UserId int64
type UserIds []UserId

type EventId int64
type EventIds []EventId
