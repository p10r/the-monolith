package domain

import (
	"log"
	"pedro-go/ra"
	"slices"
)

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

func (r *ArtistRegistry) Add(slug ra.Slug) {
	if slices.Contains(r.All().RASlugs(), slug) {
		return
	}

	res, err := r.RA.GetArtistBySlug(slug)
	if err != nil {
		log.Fatal("Oh no, handle me")
	}

	artist := Artist{RAId: res.RAID, RASlug: slug, Name: res.Name}
	r.Repo.Add(artist)
}
