package domain

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"
)

var ErrNotFoundOnRA = errors.New("artist not found on ra.com")

type ArtistRegistry struct {
	Repo    ArtistRepository
	RA      ResidentAdvisor
	Monitor EventMonitor
	Now     func() time.Time
}

func NewArtistRegistry(repo ArtistRepository, ra ResidentAdvisor, monitor EventMonitor, now func() time.Time) *ArtistRegistry {
	return &ArtistRegistry{repo, ra, monitor, now}
}

func (r *ArtistRegistry) All(ctx context.Context) Artists {
	all, err := r.Repo.All(ctx)
	if err != nil {
		log.Fatalf("error when trying to read from the db %v\n", err)
	}
	return all
}

func (r *ArtistRegistry) Follow(ctx context.Context, slug RASlug, userId UserID) error {
	all := r.All(ctx)

	i := slices.Index(all.RASlugs(), slug)
	if i != -1 {
		existing := all[i-1]
		r.Repo.Save(ctx, existing.AddFollower(userId))
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

	artist, err = r.Repo.Save(ctx, artist)
	if err != nil {
		return err
	}

	r.Monitor.Monitor(ctx, NewArtistFollowedEvent(slug, userId, r.Now))

	return nil
}

func (r *ArtistRegistry) ArtistsFor(ctx context.Context, userId UserID) (Artists, error) {
	return r.All(ctx).FilterByUserId(userId), nil
}

func (r *ArtistRegistry) AllEventsForArtist(_ context.Context, artist Artist) (Events, error) {
	now := time.Now()
	//TODO wrap error
	return r.RA.GetEventsByArtistId(artist.RAID, now, now.Add(31*24*time.Hour))
}

func (r *ArtistRegistry) NewEventsForUser(ctx context.Context, user UserID) (Events, error) {
	artists, _ := r.ArtistsFor(ctx, user)

	//TODO goroutine
	var eventsPerArtist []Events
	for _, artist := range artists {
		events, err := r.AllEventsForArtist(ctx, artist)
		if err != nil {
			return nil, fmt.Errorf("can't fetch events right now: %v", err)
		}

		for _, event := range events {
			r.Monitor.Monitor(ctx, NewNewEventForArtist(event, artist, user, r.Now))
		}

		eventsPerArtist = append(eventsPerArtist, filterAlreadyTrackedEvents(artist, events))

		_, err = r.updateEventsInDB(ctx, artist, events)
		if err != nil {
			return Events{}, err
		}
	}

	events := flatten(eventsPerArtist)

	return events, nil
}

func flatten(eventsPerArtist []Events) Events {
	var flattened Events
	for _, e := range eventsPerArtist {
		flattened = append(flattened, e...)
	}
	return flattened
}

func (r *ArtistRegistry) updateEventsInDB(ctx context.Context, artist Artist, events Events) (Events, error) {
	artist.TrackedEvents = eventIDsOf(events)
	_, err := r.Repo.Save(ctx, artist)
	if err != nil {
		return nil, fmt.Errorf("can't update events in db: %v", err)
	}
	return nil, nil
}

// TODO map ra.Event to domain type
func filterAlreadyTrackedEvents(artist Artist, events Events) Events {
	var filtered Events
	for _, e := range events {
		eventID, err := NewEventID(e.Id)
		if err != nil {
			log.Printf("failed parsing %v to EventID: %v", eventID, err)
			continue
		}
		if !artist.TrackedEvents.Contains(eventID) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func eventIDsOf(events Events) EventIDs {
	var ids EventIDs
	for _, e := range events {
		eventID, err := NewEventID(e.Id)
		if err != nil {
			log.Printf("failed parsing %v to EventID: %v", eventID, err)
			continue
		}
		ids = append(ids, eventID)
	}
	return ids
}
