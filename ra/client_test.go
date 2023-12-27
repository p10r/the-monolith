package ra_test

import (
	"bytes"
	"io"
	"net/http"
	"pedro-go/domain"
	"pedro-go/domain/expect"
	"pedro-go/ra"
	"testing"
)

func TestRAClient(t *testing.T) {
	t.Run("verify contract for in-memory fake", func(t *testing.T) {
		domain.RAContract{NewRA: func() domain.ResidentAdvisor {
			return ra.NewInMemoryClient(
				map[ra.Slug]ra.ArtistWithEvents{
					"boysnoize": {
						Artist: ra.Artist{RAID: "943", Name: "Boys Noize"},
						EventsData: []ra.Events{
							{
								Id:         "1",
								Title:      "Klubnacht",
								Date:       "2023-11-04T00:00:00.000",
								StartTime:  "2023-11-04T13:00:00.000",
								ContentUrl: "/events/1789025",
							},
							{
								Id:         "2",
								Title:      "Klubnacht 2",
								Date:       "2023-11-04T00:00:00.000",
								StartTime:  "2023-11-04T13:00:00.000",
								ContentUrl: "/events/1789025",
							},
						},
					},
					"sinamin": {
						Artist: ra.Artist{RAID: "106972", Name: "Sinamin"},
						EventsData: []ra.Events{
							{
								Id:         "1",
								Title:      "Klubnacht",
								Date:       "2023-11-04T00:00:00.000",
								StartTime:  "2023-11-04T13:00:00.000",
								ContentUrl: "/events/1789025",
							},
							{
								Id:         "2",
								Title:      "Klubnacht 2",
								Date:       "2023-11-04T00:00:00.000",
								StartTime:  "2023-11-04T13:00:00.000",
								ContentUrl: "/events/1789025",
							},
						},
					},
				},
			)
		},
		}.Test(t)
	})

	t.Run("verify contract for prod ra.co", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		domain.RAContract{NewRA: func() domain.ResidentAdvisor {
			return ra.NewClient("https://ra.co")
		}}.Test(t)
	})

	t.Run("gets artist from resident advisor", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		want := ra.Artist{RAID: "943", Name: "Boys Noize"}

		client := ra.NewClient("https://ra.co")
		got, err := client.GetArtistBySlug("boysnoize")

		expect.NoErr(t, err)
		expect.Equal(t, got, want)
	})

	t.Run("deserialize artist response", func(t *testing.T) {
		want := ra.Artist{RAID: "943", Name: "Boys Noize"}
		body := `
			{
			    "data": {
			        "artist": {
			            "id": "943",
			            "name": "Boys Noize"
			        }
			    }
			}`

		res := http.Response{Body: io.NopCloser(bytes.NewBufferString(body))}

		got, err := ra.NewArtistFrom(res.Body)

		expect.NoErr(t, err)
		expect.Equal(t, got, want)
	})
}
