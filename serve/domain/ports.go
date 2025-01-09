package domain

import (
	"context"
	"time"
)

type Flashscore interface {
	GetUpcomingMatches() (Matches, error)
}

type Discord interface {
	SendUpcomingMatches(context.Context, Matches, time.Time) error
	SendFinishedMatches(context.Context, MatchesByLeague, time.Time) error
}

type Statistics interface {
	GetItalianMenStats() StatSheets
	GetPolishMenStats() StatSheets
}
