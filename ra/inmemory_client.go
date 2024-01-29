package ra

import (
	"log"
	"pedro-go/domain"
	"time"
)

type ArtistWithEvents struct {
	Artist     Artist
	EventsData Events
}

type ArtistStore map[domain.RASlug]ArtistWithEvents

type InMemoryClient struct {
	artists ArtistStore
}

func NewInMemoryClient(artists map[domain.RASlug]ArtistWithEvents) *InMemoryClient {
	return &InMemoryClient{artists: artists}
}

func (c *InMemoryClient) GetArtistBySlug(slug domain.RASlug) (domain.ArtistInfo, error) {
	res, ok := c.artists[slug]

	if ok {
		return res.Artist.ToArtistInfo(), nil
	}

	return domain.ArtistInfo{}, ErrSlugNotFound
}

func (c *InMemoryClient) GetEventsByArtistId(
	raId string,
	_ time.Time, //TODO filter for time
	_ time.Time,
) (domain.Events, error) {
	var fil []ArtistWithEvents
	for _, a := range c.artists {
		if a.Artist.RAID == raId {
			fil = append(fil, a)
		}
	}

	if len(fil) == 0 {
		log.Fatalf("No artist found for ID %v", raId)
	}

	if len(fil) > 1 {
		log.Fatalf("More than one artist found for ID %v", raId)
	}

	first := fil[0]
	return first.EventsData.ToDomainEvents(), nil
}
