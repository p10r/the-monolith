package domain

import (
	"errors"
	"log"
	"pedro-go/ra"
	"slices"
	"time"
)

var ErrNotFoundOnRA = errors.New("artist not found on ra.com")

type ArtistRegistry struct {
	Repo ArtistRepository
	RA   ResidentAdvisor
}

func NewArtistRegistry(repo ArtistRepository, ra ResidentAdvisor) *ArtistRegistry {
	return &ArtistRegistry{repo, ra}
}

func (r *ArtistRegistry) All() Artists {
	all, err := r.Repo.All()
	if err != nil {
		log.Fatalf("error when trying to read from the db %v\n", err)
	}
	return all
}

func (r *ArtistRegistry) Follow(slug ra.Slug, userId UserId) error {
	all := r.All()
	i := slices.Index(all.RASlugs(), slug)
	if i != -1 {
		existing := all[i-1]
		r.Repo.Save(existing.AddFollower(userId))
		return nil
	}

	res, err := r.RA.GetArtistBySlug(slug)
	if err != nil {
		return ErrNotFoundOnRA
	}

	artist := Artist{
		RAId:       res.RAID,
		RASlug:     slug,
		Name:       res.Name,
		FollowedBy: UserIds{userId},
	}
	r.Repo.Save(artist)

	return nil
}

func (r *ArtistRegistry) ArtistsFor(userId UserId) (Artists, error) {
	return r.All().FilterByUserId(userId), nil
}

func (r *ArtistRegistry) EventsFor(artist Artist) ([]ra.Event, error) {
	now := time.Now()
	//TODO wrap error
	return r.RA.GetEventsByArtistId(artist.RAId, now, now.Add(9*24*time.Hour))
}
