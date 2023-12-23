package domain

import "pedro-go/ra"

type Artist struct {
	Id     int64
	RAId   string
	RASlug ra.Slug
	Name   string
}

type Artists []Artist

func (a Artists) RASlugs() []ra.Slug {
	slugs := make([]ra.Slug, len(a))
	for _, artist := range a {
		slugs = append(slugs, artist.RASlug)
	}
	return slugs
}
