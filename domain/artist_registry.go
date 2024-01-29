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

func NewDBError(err error) error {
	return fmt.Errorf("err when saving to db: %v", err)
}

func (r *ArtistRegistry) All(ctx context.Context) (Artists, error) {
	all, err := r.Repo.All(ctx)
	if err != nil {
		return nil, NewDBError(err)
	}
	return all, nil
}

func (r *ArtistRegistry) Follow(ctx context.Context, slug RASlug, userId UserID) error {
	all, err := r.All(ctx)
	if err != nil {
		return err
	}

	i := slices.Index(all.RASlugs(), slug)
	if i != -1 {
		existing := all[i-1]
		_, err := r.Repo.Save(ctx, existing.AddFollower(userId))
		if err != nil {
			return NewDBError(err)
		}
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
		return NewDBError(err)
	}

	r.Monitor.Monitor(ctx, NewArtistFollowedEvent(slug, userId, r.Now))

	return nil
}

func (r *ArtistRegistry) ArtistsFor(ctx context.Context, userId UserID) (Artists, error) {
	all, err := r.All(ctx)
	if err != nil {
		return nil, err
	}
	return all.FilterByUserId(userId), nil
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

		eventsPerArtist = append(eventsPerArtist, filterAlreadyTrackedEvents(artist, events))

		_, err = r.updateEventsInDB(ctx, artist, events)
		if err != nil {
			return Events{}, NewDBError(err)
		}

		for _, event := range events {
			r.Monitor.Monitor(ctx, NewNewEventForArtist(event, artist, user, r.Now))
		}
	}

	events := flatten(eventsPerArtist)

	return events, nil
}

func (r *ArtistRegistry) updateEventsInDB(ctx context.Context, artist Artist, events Events) (Events, error) {
	artist.TrackedEvents = eventIDsOf(events)
	_, err := r.Repo.Save(ctx, artist)
	if err != nil {
		return nil, NewDBError(err)
	}
	return nil, nil
}

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

func flatten(eventsPerArtist []Events) Events {
	var flattened Events
	for _, e := range eventsPerArtist {
		flattened = append(flattened, e...)
	}
	return flattened
}
