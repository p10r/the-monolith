package domain

import (
	"errors"
	"fmt"
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

func (r *ArtistRegistry) Follow(slug ra.Slug, userId UserID) error {
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
		RAID:       res.RAID,
		RASlug:     slug,
		Name:       res.Name,
		FollowedBy: UserIDs{userId},
	}
	r.Repo.Save(artist)

	return nil
}

func (r *ArtistRegistry) ArtistsFor(userId UserID) (Artists, error) {
	return r.All().FilterByUserId(userId), nil
}

func (r *ArtistRegistry) AllEventsForArtist(artist Artist) ([]ra.Event, error) {
	now := time.Now()
	//TODO wrap error
	return r.RA.GetEventsByArtistId(artist.RAID, now, now.Add(31*24*time.Hour))
}

func (r *ArtistRegistry) NewEventsForUser(id UserID) ([]ra.Event, error) {
	artists, _ := r.ArtistsFor(id)

	//TODO goroutine
	var eventsPerArtist [][]ra.Event
	for _, artist := range artists {
		e, err := r.AllEventsForArtist(artist)
		if err != nil {
			return nil, fmt.Errorf("can't fetch events right now: %v", err)
		}
		eventsPerArtist = append(eventsPerArtist, e)
	}

	var flattened []ra.Event
	for _, e := range eventsPerArtist {
		flattened = append(flattened, e...)
	}

	return flattened, nil
}
