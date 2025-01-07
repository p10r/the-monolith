package domain

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
	StatsUrl         string // Set by statistics package
}

func (m Match) LeagueKey() LeagueKey {
	return LeagueKey(strings.ToLower(m.FlashscoreName))
}

type Matches []Match

func (matches Matches) Scheduled() Matches {
	scheduled := Matches{}

	for _, match := range matches {
		if lowerCase(match.Stage) == lowerCase("SCHEDULED") {
			scheduled = append(scheduled, match)
		}
	}
	if len(scheduled) == 0 {
		return Matches{}
	}
	return scheduled
}

func (matches Matches) Finished() Matches {
	finished := Matches{}
	for _, match := range matches {
		if lowerCase(match.Stage) == lowerCase("FINISHED") {
			finished = append(finished, match)
		}
	}

	if len(finished) == 0 {
		return Matches{}
	}

	return finished
}

func lowerCase(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

type MatchesByLeague map[LeagueKey]Matches

func (matches Matches) ToMap() MatchesByLeague {
	out := map[LeagueKey]Matches{}
	for _, match := range matches {
		out[match.LeagueKey()] = append(out[match.LeagueKey()], match)
	}
	return out
}
