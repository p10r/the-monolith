package serve

import (
	"strings"
)

type LeagueKey string

func (k LeagueKey) CountryEquals(country string) bool {
	keyCountry := strings.Split(string(k), ":")[0]

	return strings.EqualFold(keyCountry, country)
}

func NewLeagueKey(fullFlashscoreName string) LeagueKey {
	return LeagueKey(fullFlashscoreName)
}

type Match struct {
	HomeName         string
	AwayName         string
	StartTime        int64
	FlashscoreName   string // Country + League
	Country          string
	League           string
	Stage            string
	HomeScoreCurrent int
	AwayScoreCurrent int
}

func (m Match) LeagueKey() LeagueKey {
	return LeagueKey(strings.ToLower(m.FlashscoreName))
}

type Matches []Match

func (matches Matches) Scheduled() Matches {
	return matches.withStage("SCHEDULED")
}

func (matches Matches) Finished() Matches {
	return matches.withStage("FINISHED")
}

func (matches Matches) withStage(stage string) Matches {
	var filtered Matches
	for _, match := range matches {
		if strings.EqualFold(match.Stage, stage) {
			filtered = append(filtered, match)
		}
	}
	return filtered
}

type MatchesByLeague map[LeagueKey]Matches

func (matches Matches) ToMap() MatchesByLeague {
	out := map[LeagueKey]Matches{}
	for _, match := range matches {
		out[match.LeagueKey()] = append(out[match.LeagueKey()], match)
	}
	return out
}
