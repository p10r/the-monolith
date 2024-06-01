package serve

import (
	"context"
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/pkg/sqlite"
	"github.com/p10r/pedro/serve/db"
	"github.com/p10r/pedro/serve/discord"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/flashscore"
	"github.com/robfig/cron/v3"
	"log/slog"
	"time"
)

type Serve struct {
	importer       *domain.MatchImporter
	importSchedule string // CRON
	log            *slog.Logger
}

// NewServeApp wires Serve App together.
// Expects an already opened connection.
func NewServeApp(
	conn *sqlite.DB,
	flashscoreUri, flashscoreApiKey, discordUri string,
	favouriteLeagues []string,
	logHandler slog.Handler,
	importSchedule string,
) *Serve {
	log := l.NewAppLogger(logHandler, "serve")

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

	store := db.NewMatchStore(conn)
	flashscoreClient := flashscore.NewClient(flashscoreUri, flashscoreApiKey, log)
	discordClient := discord.NewClient(discordUri, log)
	now := func() time.Time { return time.Now() }

	importer := domain.NewMatchImporter(
		store,
		flashscoreClient,
		discordClient,
		favouriteLeagues,
		now,
		log,
	)

	return &Serve{importer, importSchedule, log}
}

func (s *Serve) Start(ctx context.Context) {
	c := cron.New()
	_, err := c.AddFunc(s.importSchedule, func() {
		_, err := s.importer.ImportScheduledMatches(ctx)
		if err != nil {
			s.log.Error(l.Error("serve run failed", err))
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
