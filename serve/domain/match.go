package domain

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

type Matches []Match

type FinishedMatch struct {
	Match
	StatsUrl string
}

type FinishedMatches []FinishedMatch
