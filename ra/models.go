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

type Event struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Date       string `json:"date"`
	StartTime  string `json:"startTime"`
	ContentUrl string `json:"contentUrl"`
	Images     []struct {
		Filename string `json:"filename"`
	} `json:"images"`
	Venue struct {
		Id         string `json:"id"`
		Name       string `json:"name"`
		ContentUrl string `json:"contentUrl"`
		Area       struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"area"`
	} `json:"venue"`
}

type ErrorRes struct {
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Extensions struct {
			Code string `json:"code"`
		} `json:"extensions"`
	} `json:"errors"`
}
