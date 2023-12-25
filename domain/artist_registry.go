package domain

import (
	"errors"
	"fmt"
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

func (r *ArtistRegistry) Follow(slug ra.Slug, userId UserId) error {
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
		RAId:       res.RAID,
		RASlug:     slug,
		Name:       res.Name,
		FollowedBy: UserIds{userId},
	}
	r.Repo.Save(artist)

	return nil
}

func (r *ArtistRegistry) ArtistsFor(userId UserId) (Artists, error) {
	all := r.All()
	fmt.Printf("artists for: %v\n", all)
	return all.FilterByUserId(userId), nil
}
