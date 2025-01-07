package domain_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/testutil"
	"testing"
)

func TestDomain(t *testing.T) {
	t.Run("filters for scheduled matches", func(t *testing.T) {
		input := domain.Matches{
			{
				HomeName:       "Trentino",
				AwayName:       "Jastrzebski",
				FlashscoreName: "Europe: Champions League - Play Offs",
				Stage:          "SCHEDULED",
			},
			{
				HomeName:       "Resovia",
				AwayName:       "Zaksa",
				FlashscoreName: "Europe: Champions League - Play Offs",
				Stage:          "FINISHED",
			},
		}

		expected := domain.Matches{
			{
				HomeName:       "Trentino",
				AwayName:       "Jastrzebski",
				FlashscoreName: "Europe: Champions League - Play Offs",
				Stage:          "SCHEDULED",
			},
		}

		matches := input.Scheduled()
		assert.Equal(t, matches, expected)
	})

	t.Run("filters for finished matches", func(t *testing.T) {
		input := domain.Matches{
			{
				HomeName:       "Trentino",
				AwayName:       "Jastrzebski",
				FlashscoreName: "Europe: Champions League - Play Offs",
				Stage:          "SCHEDULED",
			},
			{
				HomeName:       "Resovia",
				AwayName:       "Zaksa",
				FlashscoreName: "Europe: Champions League - Play Offs",
				Stage:          "FINISHED",
			},
		}

		expected := domain.Matches{
			{
				HomeName:       "Resovia",
				AwayName:       "Zaksa",
				FlashscoreName: "Europe: Champions League - Play Offs",
				Stage:          "FINISHED",
			},
		}

		matches := input.Finished()
		assert.Equal(t, matches, expected)
	})

	t.Run("handles 0 scheduled matches", func(t *testing.T) {
		m := domain.Matches{}.Scheduled()
		assert.Equal(t, len(m), 0)
	})

	t.Run("filters for favourites", func(t *testing.T) {
		expected := domain.Matches{
			{
				HomeName:         "Trentino",
				AwayName:         "Jastrzebski",
				StartTime:        1714917600,
				FlashscoreName:   "Europe: Champions League - Play Offs",
				Country:          "Europe",
				League:           "Champions League - Play Offs",
				Stage:            "SCHEDULED",
				HomeScoreCurrent: 3,
				AwayScoreCurrent: 0,
			},
			{
				HomeName:         "Resovia",
				AwayName:         "Zaksa",
				StartTime:        1714917600,
				FlashscoreName:   "Europe: Champions League - Play Offs",
				Country:          "Europe",
				League:           "Champions League - Play Offs",
				Stage:            "SCHEDULED",
				HomeScoreCurrent: 3,
				AwayScoreCurrent: 0,
			},
			{
				HomeName:         "Grand Rapids Rise W",
				AwayName:         "San Diego Mojo W",
				StartTime:        1714939200,
				FlashscoreName:   "USA: PVF Women",
				Country:          "USA",
				League:           "PVF Women",
				Stage:            "SCHEDULED",
				HomeScoreCurrent: 3,
				AwayScoreCurrent: 0,
			},
		}

		matches := testutil.Matches(t).Scheduled()
		assert.Equal(t, matches, expected)
	})
}
