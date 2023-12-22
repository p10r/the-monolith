package ra

import (
	"encoding/json"
	"io"
)

type ArtistResWrapper struct {
	ArtistData `json:"data"`
}

type ArtistData struct {
	ArtistRes `json:"artist"`
}

type ArtistRes struct {
	RAID string `json:"id"`
	Name string `json:"name"`
}

func NewArtistFrom(input io.ReadCloser) (ArtistRes, error) {
	var res ArtistResWrapper
	err := json.NewDecoder(input).Decode(&res)

	return res.ArtistRes, err
}
