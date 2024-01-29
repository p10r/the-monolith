package domain

import (
	"context"
	"encoding/json"
	"time"
)

// EventMonitor tracks application events. Doesn't return errors , but logs them.
type EventMonitor interface {
	Monitor(ctx context.Context, e MonitoringEvent)
	All(ctx context.Context) (MonitoringEvents, error)
}

type MonitoringEvent interface {
	Name() string
	ToJSON() ([]byte, error)
	Timestamp() time.Time
}

type MonitoringEvents []MonitoringEvent

type ArtistFollowed struct {
	ArtistSlug string
	UserId     UserID
	CreatedAt  time.Time
}

func NewArtistFollowedEvent(slug RASlug, id UserID, now func() time.Time) ArtistFollowed {
	return ArtistFollowed{string(slug), id, now()}
}

func (e ArtistFollowed) Name() string {
	return "ArtistFollowed"
}

func (e ArtistFollowed) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e ArtistFollowed) Timestamp() time.Time {
	return e.CreatedAt
}

type NewEventForArtist struct {
	EventId   string
	Slug      string
	User      UserID
	CreatedAt time.Time
}

func NewNewEventForArtist(
	event Event,
	artist Artist,
	id UserID,
	now func() time.Time,
) NewEventForArtist {
	return NewEventForArtist{event.Id, string(artist.RASlug), id, now()}
}

func (e NewEventForArtist) Name() string {
	return "NewEventForArtist"
}

func (e NewEventForArtist) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e NewEventForArtist) Timestamp() time.Time {
	return e.CreatedAt
}
