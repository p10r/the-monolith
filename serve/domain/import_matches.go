package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/p10r/pedro/pkg/l"
	"log/slog"
	"time"
)

type Flashscore interface {
	GetUpcomingMatches() (UntrackedMatches, error)
}

type Discord interface {
	SendUpcomingMatches(context.Context, Matches, time.Time) error
	SendFinishedMatches(context.Context, FinishedMatches, time.Time) error
}

type MatchImporter struct {
	store      MatchStore
	flashscore Flashscore
	discord    Discord
	favLeagues []string
	clock      func() time.Time
	log        *slog.Logger
}

func NewMatchImporter(
	store MatchStore,
	flashscore Flashscore,
	discord Discord,
	favLeagues []string,
	clock func() time.Time,
	log *slog.Logger,
) *MatchImporter {
	return &MatchImporter{
		store,
		flashscore,
		discord,
		favLeagues,
		clock,
		log,
	}
}

// ImportScheduledMatches writes matches from flashscore into the db for the current day.
// Doesn't validate if the match is already present,
// as it's expected to be triggered only once per day for now.
func (importer *MatchImporter) ImportScheduledMatches(ctx context.Context) (Matches, error) {
	untrackedMatches, err := importer.fetchAllMatches()
	if err != nil {
		importer.log.Error(l.Error("cannot fetch matches", err))
		return nil, err
	}

	//TODO remove error, return empty slice
	upcoming := untrackedMatches.FilterScheduled(importer.favLeagues)
	if len(upcoming) == 0 {
		importer.log.Info("No upcoming games today")
		return Matches{}, nil
	}

	importer.log.Info(fmt.Sprintf("%v matches coming up today", len(upcoming)))

	trackedMatches, err := importer.storeUntrackedMatches(ctx, upcoming)
	if err != nil {
		importer.log.Error(l.Error("error when writing to db", err))
		return nil, err
	}

	err = importer.discord.SendUpcomingMatches(ctx, trackedMatches, importer.clock())
	if err != nil {
		importer.log.Error(l.Error("send to discord error", err))
		return nil, err
	}

	return trackedMatches, nil
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

	err = importer.discord.SendFinishedMatches(ctx, finished, importer.clock())
	if err != nil {
		importer.log.Error(l.Error("send to discord error", err))
		return err
	}

	return nil
}

func (importer *MatchImporter) fetchAllMatches() (UntrackedMatches, error) {
	untrackedMatches, err := importer.flashscore.GetUpcomingMatches()
	if err != nil {
		return nil, fmt.Errorf("could not fetch matches from flashscore: %v", err)
	}
	return untrackedMatches, err
}

func (importer *MatchImporter) storeUntrackedMatches(
	ctx context.Context,
	matches UntrackedMatches,
) (Matches, error) {
	var trackedMatches Matches
	var dbErrs []error
	for _, untrackedMatch := range matches {
		trackedMatch, err := importer.store.Add(ctx, untrackedMatch)
		if err != nil {
			//nolint
			dbErr := fmt.Errorf("could not persist match %v, aborting: %v", untrackedMatch.HomeName, err)
			dbErrs = append(dbErrs, dbErr)
		}

		importer.log.Debug(fmt.Sprintf("Stored in DB: %v", trackedMatch))

		trackedMatches = append(trackedMatches, trackedMatch)
	}

	return trackedMatches, errors.Join(dbErrs...)
}
