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
	t.Run("test contract against real ra.co", func(t *testing.T) {
		domain.RAContract{NewRA: func() domain.ResidentAdvisor {
			return ra.NewInMemoryClient(
				map[ra.Slug]ra.Artist{
					"boysnoize": {RAID: "943", Name: "Boys Noize"},
				},
			)
		}}.Test(t)
	})

	t.Run("test contract against real ra.co", func(t *testing.T) {
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
