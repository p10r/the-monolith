package domain

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestStatistics(t *testing.T) {
	for _, tt := range []struct {
		tc           string
		leagueFilter []string
		matches      Matches
		stats        StatSheets
		expected     Matches
	}{
		{
			tc:           "zips same matches",
			leagueFilter: []string{"poland"},
			matches: Matches{
				Match{
					HomeName:       "Zawierce",
					AwayName:       "Warsaw",
					FlashscoreName: "poland",
				},
			},
			stats: StatSheets{
				StatSheet{
					Home: "Zawierce",
					Away: "Warsaw",
					Url:  "/a-status-url",
				},
			},
			expected: Matches{
				Match{
					HomeName:       "Zawierce",
					AwayName:       "Warsaw",
					FlashscoreName: "poland",
					StatsUrl:       "/a-status-url",
				},
			},
		},

		{
			tc:           "zips only matching matches",
			leagueFilter: []string{"poland"},
			matches: Matches{
				Match{
					HomeName:       "Zawierce",
					AwayName:       "Warsaw",
					FlashscoreName: "poland",
				},
				Match{
					HomeName:       "Something Else",
					AwayName:       "Another one",
					FlashscoreName: "poland",
				},
			},
			stats: StatSheets{
				StatSheet{
					Home: "Zawierce",
					Away: "Warsaw",
					Url:  "/a-status-url",
				},
			},
			expected: Matches{
				Match{
					HomeName:       "Zawierce",
					AwayName:       "Warsaw",
					FlashscoreName: "poland",
					StatsUrl:       "/a-status-url",
				},
				Match{
					HomeName:       "Something Else",
					AwayName:       "Another one",
					FlashscoreName: "poland",
				},
			},
		},
	} {
		t.Run(tt.tc, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.matches.ZipWith(tt.stats))
		})
	}
}
