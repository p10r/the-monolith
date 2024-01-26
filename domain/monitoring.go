package domain

import (
	"context"
	"encoding/json"
)

type EventMonitor interface {
	Monitor(ctx context.Context, e MonitoringEvent)
	All(ctx context.Context) (MonitoringEvents, error)
}

type MonitoringEvent interface {
	Name() string
	ToJSON() ([]byte, error)
}

type MonitoringEvents []MonitoringEvent

type ArtistFollowedEvent struct {
	ArtistSlug string
	UserId     UserID
}

func NewArtistFollowedEvent(slug RASlug, id UserID) ArtistFollowedEvent {
	return ArtistFollowedEvent{string(slug), id}
}

func (e ArtistFollowedEvent) Name() string {
	return "ArtistFollowedEvent"
}

func (e ArtistFollowedEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

type NewEventForArtistEvent struct {
	EventId string
	Slug    string
	Users   UserIDs
}

func NewNewEventForArtistEvent(event Event, artist Artist, ids UserIDs) NewEventForArtistEvent {
	return NewEventForArtistEvent{event.Id, string(artist.RASlug), ids}
}

func (e NewEventForArtistEvent) Name() string {
	return "NewEventForArtistEvent"
}

func (e NewEventForArtistEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
