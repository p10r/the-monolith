package domain

import (
	"context"
	"fmt"
	"github.com/p10r/pedro/pkg/l"
	"log/slog"
	"time"
)

type Flashscore interface {
	GetUpcomingMatches() (Matches, error)
}

type Discord interface {
	SendUpcomingMatches(context.Context, Matches, time.Time) error
	SendFinishedMatches(context.Context, FinishedMatchesByLeague, time.Time) error
}

type Statistics interface {
	EnrichMatches(matches FinishedMatchesByLeague) FinishedMatchesByLeague
}

type MatchImporter struct {
	flashscore Flashscore
	discord    Discord
	statistics Statistics
	favLeagues []string
	clock      func() time.Time
	log        *slog.Logger
}

func NewMatchImporter(
	flashscore Flashscore,
	discord Discord,
	statistics Statistics,
	favLeagues []string,
	clock func() time.Time,
	log *slog.Logger,
) *MatchImporter {
	return &MatchImporter{
		flashscore,
		discord,
		statistics,
		favLeagues,
		clock,
		log,
	}
}

// ImportScheduledMatches writes matches from flashscore into the db for the current day.
// Doesn't validate if the match is already present,
// as it's expected to be triggered only once per day for now.
func (importer *MatchImporter) ImportScheduledMatches(ctx context.Context) (Matches, error) {
	matches, err := importer.fetchAllMatches()
	if err != nil {
		importer.log.Error(l.Error("cannot fetch matches", err))
		return nil, err
	}

	//TODO remove error, return empty slice
	upcoming := matches.FilterScheduled(importer.favLeagues)
	if len(upcoming) == 0 {
		importer.log.Info("No upcoming games today")
		return Matches{}, nil
	}

	err = importer.discord.SendUpcomingMatches(ctx, upcoming, importer.clock())
	if err != nil {
		importer.log.Error(l.Error("send to discord error", err))
		return nil, err
	}

	return upcoming, nil
}

func (importer *MatchImporter) ImportFinishedMatches(ctx context.Context) error {
	flashscoreMatches, err := importer.fetchAllMatches()
	if err != nil {
		importer.log.Error(l.Error("cannot fetch matches", err))
		return err
	}

	finished := flashscoreMatches.FilterFinished(importer.favLeagues)
	if len(finished) == 0 {
		importer.log.Info("No finished games today")
		return nil
	}

	matchesWithStats := importer.statistics.EnrichMatches(finished.ToMap())

	err = importer.discord.SendFinishedMatches(ctx, matchesWithStats, importer.clock())
	if err != nil {
		importer.log.Error(l.Error("send to discord error", err))
		return err
	}

	return nil
}

func (importer *MatchImporter) fetchAllMatches() (Matches, error) {
	m, err := importer.flashscore.GetUpcomingMatches()
	if err != nil {
		return nil, fmt.Errorf("could not fetch matches from flashscore: err: %v", err)
	}
	return m, err
}
