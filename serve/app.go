package serve

import (
	"context"
	"github.com/p10r/pedro/pkg/httputil"
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/serve/discord"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/flashscore"
	"github.com/p10r/pedro/serve/statistics"
	"github.com/robfig/cron/v3"
	"log/slog"
	"time"
)

type Serve struct {
	importer               *domain.MatchImporter
	importUpcomingSchedule string // CRON
	importFinishedSchedule string
	log                    *slog.Logger
}

func NewServeProdApp(
	flashscoreUri, flashscoreApiKey, discordUri string,
	log *slog.Logger,
	importUpcomingSchedule string,
	importFinishedSchedule string,
) *Serve {
	return NewServeApp(
		flashscoreUri,
		flashscoreApiKey,
		discordUri,
		log,
		importUpcomingSchedule,
		importFinishedSchedule,
	)
}

// NewServeApp wires Serve App together.
// Expects an already opened connection.
func NewServeApp(
	flashscoreUri, flashscoreApiKey, discordUri string,
	log *slog.Logger,
	importUpcomingSchedule string,
	importFinishedSchedule string,
) *Serve {
	log.Info("Starting Serve App")

	if flashscoreUri == "" {
		log.Error("flashscoreUri has not been set")
	}
	if flashscoreApiKey == "" {
		log.Error("flashscoreApiKey has not been set")
	}
	if discordUri == "" {
		log.Error("DISCORD_URI has not been set")
	}

	flashscoreClient := flashscore.NewClient(flashscoreUri, flashscoreApiKey, log)
	discordClient := discord.NewClient(discordUri, log)
	stats := statistics.NewAggregator(
		"https://www.plusliga.pl",
		"https://www.legavolley.it/",
		log,
		httputil.NewDefaultClient(),
	)
	now := func() time.Time { return time.Now() }

	importer := domain.NewMatchImporter(
		flashscoreClient,
		discordClient,
		stats,
		now,
		log,
	)

	return &Serve{
		importer,
		importUpcomingSchedule,
		importFinishedSchedule,
		log,
	}
}

func (s *Serve) StartBackgroundJobs(ctx context.Context) {
	c := cron.New()
	_, err := c.AddFunc(s.importUpcomingSchedule, func() {
		_, err := s.importer.ImportScheduledMatches(ctx)
		if err != nil {
			s.log.Error(l.Error("serve import scheduled failed", err))
		}
	})
	if err != nil {
		s.log.Error(l.Error("serve run failed", err))
	}

	_, err = c.AddFunc(s.importFinishedSchedule, func() {
		err = s.importer.ImportFinishedMatches(ctx)
		if err != nil {
			s.log.Error(l.Error("serve import finished failed", err))
		}
	})
	if err != nil {
		s.log.Error(l.Error("serve run failed", err))
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.Start()
		}
	}
}
