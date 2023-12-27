package ra

import (
	"log"
	"time"
)

type ArtistWithEvents struct {
	Artist     Artist
	EventsData []Event
}

type InMemoryClient struct {
	artists map[Slug]ArtistWithEvents
}

func NewInMemoryClient(artists map[Slug]ArtistWithEvents) *InMemoryClient {
	return &InMemoryClient{artists: artists}
}

func (c *InMemoryClient) GetArtistBySlug(slug Slug) (Artist, error) {
	res, ok := c.artists[slug]

	if ok {
		return res.Artist, nil
	}

	return Artist{}, ErrSlugNotFound
}

func (c *InMemoryClient) GetEventsByArtistId(
	raId string,
	_ time.Time,
	_ time.Time,
) ([]Event, error) {
	var fil []ArtistWithEvents
	for _, a := range c.artists {
		if a.Artist.RAID == raId {
			fil = append(fil, a)
		}
	}

	if len(fil) == 0 {
		log.Fatalf("No artist found for Id %v", raId)
	}

	if len(fil) > 1 {
		log.Fatalf("More than one artist found for Id %v", raId)
	}

	return fil[0].EventsData, nil
}
