package ra

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"pedro-go/domain"
)

type Artist struct {
	RAID string `json:"id"`
	Name string `json:"name"`
}

func NewArtist(res *http.Response) (Artist, error) {
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Artist{}, fmt.Errorf("request failed with status code: %v", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Artist{}, fmt.Errorf("cannot parse response body")
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
