package ra

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Artist struct {
	RAID string `json:"id"`
	Name string `json:"name"`
}

func NewArtist(res *http.Response, err error) (Artist, error) {
	if err != nil {
		return Artist{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Artist{}, fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	var body struct {
		Data struct {
			Artist `json:"artist"`
		} `json:"data"`
	}
	err = json.NewDecoder(res.Body).Decode(&body)
	if body.Data.Artist == (Artist{}) {
		return Artist{}, ErrSlugNotFound
	}

	return body.Data.Artist, err
}
