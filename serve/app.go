package serve

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/p10r/monolith/pkg/l"
	"github.com/robfig/cron/v3"
)

// Serve fetches volleyball matches from Flashscore and posts them to Discord.
type Serve struct {
	flashscoreBaseURL string
	flashscoreAPIKey  string
	discordURL        string
	httpClient        *http.Client
	now               func() time.Time
	log               *slog.Logger

	importUpcomingSchedule string
	importFinishedSchedule string
}

func New(
	flashscoreBaseURL, flashscoreAPIKey, discordURL string,
	log *slog.Logger,
	importUpcomingSchedule, importFinishedSchedule string,
) *Serve {
	if flashscoreBaseURL == "" {
		log.Error("flashscoreUri has not been set")
	}
	if flashscoreAPIKey == "" {
		log.Error("flashscoreApiKey has not been set")
	}
	if discordURL == "" {
		log.Error("DISCORD_URI has not been set")
	}

	return &Serve{
		flashscoreBaseURL: flashscoreBaseURL,
		flashscoreAPIKey:  flashscoreAPIKey,
		discordURL:        discordURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 10 * time.Second,
				}).DialContext,
				TLSHandshakeTimeout:   5 * time.Second,
				ResponseHeaderTimeout: 5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
		now:                    time.Now,
		log:                    log.With(slog.String("app", "serve")),
		importUpcomingSchedule: importUpcomingSchedule,
		importFinishedSchedule: importFinishedSchedule,
	}
}

func (s *Serve) StartBackgroundJobs(ctx context.Context) {
	c := cron.New()

	if _, err := c.AddFunc(s.importUpcomingSchedule, func() {
		if _, err := s.importScheduledMatches(); err != nil {
			s.log.Error(l.Error("serve import scheduled failed", err))
		}
	}); err != nil {
		s.log.Error(l.Error("add serve import scheduled job", err))
	}

	if _, err := c.AddFunc(s.importFinishedSchedule, func() {
		if err := s.importFinishedMatches(); err != nil {
			s.log.Error(l.Error("serve import finished job", err))
		}
	}); err != nil {
		s.log.Error(l.Error("add import finished job", err))
	}

	c.Start()
	<-ctx.Done()
	c.Stop()
}

func (s *Serve) importScheduledMatches() (Matches, error) {
	matches, err := s.fetchMatches()
	if err != nil {
		return nil, err
	}

	scheduled := matches.Favourites().Scheduled()
	if len(scheduled) == 0 {
		s.log.Info("No scheduled matches today")
		return nil, nil
	}

	if err := s.sendScheduledMatches(scheduled, s.now()); err != nil {
		s.log.Error(l.Error("send scheduled matches to discord", err))
		return nil, err
	}

	return scheduled, nil
}

func (s *Serve) importFinishedMatches() error {
	matches, err := s.fetchMatches()
	if err != nil {
		return err
	}

	finished := matches.Favourites().Finished()
	if len(finished) == 0 {
		s.log.Info("No finished matches today")
		return nil
	}

	if err := s.sendFinishedMatches(finished.ToMap(), s.now()); err != nil {
		s.log.Error(l.Error("send finished matches to discord", err))
		return err
	}

	return nil
}
