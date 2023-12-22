package domain

type ArtistRegistry struct {
	Repo ArtistRepository
}

func NewArtistRegistry(repo ArtistRepository) *ArtistRegistry {
	return &ArtistRegistry{repo}
}

func (r *ArtistRegistry) All() []Artist {
	return r.Repo.All()
}
