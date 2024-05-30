package ra

import (
	"encoding/json"
	"fmt"
	"github.com/p10r/pedro/pedro/domain"
	"io"
	"net/http"
)

type Artist struct {
	RAID string `json:"id"`
	Name string `json:"name"`
}

func NewArtist(res *http.Response) (Artist, error) {
	if res == nil {
		return Artist{}, fmt.Errorf("artist response was null")
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Artist{}, fmt.Errorf("cannot parse response body")
	}

	if res.StatusCode != http.StatusOK {
		//nolint:lll
		return Artist{}, fmt.Errorf("request failed with status code: %v, body: %v", res.StatusCode, string(data))
	}

	var body struct {
		Data struct {
			Artist `json:"artist"`
		} `json:"data"`
	}

	err = json.Unmarshal(data, &body)
	if err != nil {
		return Artist{}, fmt.Errorf("JSON deserialization error. Body: %s", data)
	}

	if body.Data.Artist == (Artist{}) {
		return Artist{}, ErrSlugNotFound
	}

	return body.Data.Artist, err
}

func (a Artist) ToArtistInfo() domain.ArtistInfo {
	return domain.ArtistInfo{
		RAID: a.RAID,
		Name: a.Name,
	}
}
