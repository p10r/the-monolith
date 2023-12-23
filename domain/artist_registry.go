package domain

import (
	"errors"
	"pedro-go/ra"
	"slices"
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
	return r.Repo.All()
}

func (r *ArtistRegistry) Add(slug ra.Slug) error {
	if slices.Contains(r.All().RASlugs(), slug) {
		return nil
	}

	res, err := r.RA.GetArtistBySlug(slug)
	if err != nil {
		return ErrNotFoundOnRA
	}

	artist := Artist{RAId: res.RAID, RASlug: slug, Name: res.Name}
	r.Repo.Add(artist)

	return nil
}
