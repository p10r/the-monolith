package ra

import (
	"github.com/p10r/monolith/pedro/domain"
	"github.com/p10r/monolith/pkg/l"
	"testing"
	"time"
)

type ArtistWithEvents struct {
	Artist     Artist
	EventsData Events
}

type ArtistStore map[domain.RASlug]ArtistWithEvents

type InMemoryClient struct {
	artists ArtistStore
	t       *testing.T
}

func NewInMemoryClient(
	t *testing.T,
	artists map[domain.RASlug]ArtistWithEvents,
) *InMemoryClient {
	return &InMemoryClient{artists: artists, t: t}
}

func (c *InMemoryClient) GetArtistBySlug(slug domain.RASlug) (domain.ArtistInfo, error) {
	res, ok := c.artists[slug]

	if ok {
		return res.Artist.ToArtistInfo(), nil
	}

	return domain.ArtistInfo{}, ErrSlugNotFound
}

func (c *InMemoryClient) GetEventsByArtistId(
	a domain.Artist,
	_ time.Time, //TODO filter for time
	_ time.Time,
) (domain.Events, error) {
	log := l.NewTextLogger()

	raId := a.RAID
	var fil []ArtistWithEvents
	for _, a := range c.artists {
		if a.Artist.RAID == raId {
			fil = append(fil, a)
		}
	}

	if len(fil) == 0 {
		c.t.Fatalf("No artist found for ID %v", raId)
	}

	if len(fil) > 1 {
		c.t.Fatalf("More than one artist found for ID %v", raId)
	}

	first := fil[0]
	return first.EventsData.ToDomainEvents(a.Name, log), nil
}
