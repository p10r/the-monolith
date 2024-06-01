package ra_test

import (
	"github.com/p10r/pedro/pedro/domain"
	"github.com/p10r/pedro/pedro/domain/expect"
	"github.com/p10r/pedro/pedro/ra"
	"log/slog"
	"os"
	"testing"
)

func TestRAClient(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Run("verify contract for in-memory fake", func(t *testing.T) {
		domain.RAContract{NewRA: func() domain.ResidentAdvisor {
			return ra.NewInMemoryClient(t,
				map[domain.RASlug]ra.ArtistWithEvents{
					"boysnoize": {
						Artist: ra.Artist{RAID: "943", Name: "Boys Noize"},
						EventsData: []ra.Event{
							{
								Id:         "1",
								Title:      "Klubnacht",
								Date:       "2023-11-04T00:00:00.000",
								StartTime:  "2023-11-04T13:00:00.000",
								ContentUrl: "/events/1789025",
								Venue: ra.Venue{
									Area: ra.Area{
										Name: "Berlin",
									},
									Name: "RSO",
								},
							},
							{
								Id:         "2",
								Title:      "Klubnacht 2",
								Date:       "2023-11-04T00:00:00.000",
								StartTime:  "2023-11-04T13:00:00.000",
								ContentUrl: "/events/1789025",
								Venue: ra.Venue{
									Area: ra.Area{
										Name: "Berlin",
									},
									Name: "RSO",
								},
							},
						},
					},
					"sinamin": {
						Artist: ra.Artist{RAID: "106972", Name: "Sinamin"},
						EventsData: []ra.Event{
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
			return ra.NewClient("https://ra.co", log)
		}}.Test(t)
	})

	t.Run("gets artist from resident advisor", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}

		want := domain.ArtistInfo{RAID: "943", Name: "Boys Noize"}

		client := ra.NewClient("https://ra.co", log)
		got, err := client.GetArtistBySlug("boysnoize")

		expect.NoErr(t, err)
		expect.Equal(t, got, want)
	})
}
