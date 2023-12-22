package ra

import (
	"encoding/json"
	"io"
)

type Slug string

type ArtistResWrapper struct {
	ArtistData `json:"data"`
}

type ArtistData struct {
	Artist `json:"artist"`
}

type Artist struct {
	RAID string `json:"id"`
	Name string `json:"name"`
}

func NewArtistFrom(input io.ReadCloser) (Artist, error) {
	var res ArtistResWrapper
	err := json.NewDecoder(input).Decode(&res)

	return res.Artist, err
}
