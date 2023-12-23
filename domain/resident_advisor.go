package domain

import "pedro-go/ra"

type ResidentAdvisor interface {
	GetArtistBySlug(slug ra.Slug) (ra.Artist, error)
}
