package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/p10r/monolith/pkg/l"
	"github.com/p10r/monolith/serve"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	TelegramToken    string `env:"TELEGRAM_TOKEN"`
	TelegramAdmin    string `env:"TELEGRAM_ADMIN_USER"`
	FlashscoreApiKey string `env:"FLASHSCORE_API_KEY"`
	DiscordUri       string `env:"DISCORD_URI"`
}

const serveImportUpcomingSchedule = "0 6 * * *"
const serveImportFinishedSchedule = "0 22 * * *"
const flashscoreUri = "https://flashscore.p.rapidapi.com"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Fly.io adds its own timestamp to logs, so drop slog's.
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				return slog.Attr{}
			}
			return a
		},
	})
	log := slog.New(handler)

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Error(l.Error("error loading config", err))
	}

	serve.New(
		flashscoreUri,
		cfg.FlashscoreApiKey,
		cfg.DiscordUri,
		l.NewTelegramLogger(cfg.TelegramToken, cfg.TelegramAdmin, "serve"),
		serveImportUpcomingSchedule,
		serveImportFinishedSchedule,
	).StartBackgroundJobs(ctx)
}
