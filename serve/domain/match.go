package domain

import "strings"

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

type FinishedMatch struct {
	Match
	StatsUrl string
}

func (m FinishedMatch) LeagueKey() LeagueKey {
	return LeagueKey(strings.ToLower(m.FlashscoreName))
}

type FinishedMatches []FinishedMatch

type FinishedMatchesByLeague map[LeagueKey]FinishedMatches

func (m FinishedMatches) ToMap() FinishedMatchesByLeague {
	out := map[LeagueKey]FinishedMatches{}
	for _, match := range m {
		out[match.LeagueKey()] = append(out[match.LeagueKey()], match)
	}
	return out
}
