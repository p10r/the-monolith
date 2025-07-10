package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/p10r/monolith/giftbox"
	"github.com/p10r/monolith/pkg/l"
	"github.com/p10r/monolith/pkg/sqlite"
	"github.com/p10r/monolith/serve"
	gracefulshutdown "github.com/quii/go-graceful-shutdown"
	"github.com/sethvargo/go-envconfig"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Config struct {
	TelegramToken    string  `env:"TELEGRAM_TOKEN"`
	TelegramAdmin    string  `env:"TELEGRAM_ADMIN_USER"`
	DSN              string  `env:"DSN"`
	AllowedUserIds   []int64 `env:"ALLOWED_USER_IDS"`
	FlashscoreApiKey string  `env:"FLASHSCORE_API_KEY"`
	DiscordUri       string  `env:"DISCORD_URI"`
	GiftBoxApiKey    string  `env:"GIFT_BOX_API_KEY"`
	ServeApiKey      string  `env:"SERVE_API_KEY"`
}

const serveImportUpcomingSchedule = "0 6 * * *"
const serveImportFinishedSchedule = "0 22 * * *"
const flashscoreUri = "https://flashscore.p.rapidapi.com"

func main() {
	// Remove time, as fly.io is logging it automatically
	// Time is automatically left out if it is 0, see HandlerOptions.ReplaceAttr
	// For more examples, see
	// https://github.com/golang/go/blob/master/src/log/slog/example_custom_levels_test.go
	//
	// Good article for further configuration:
	// https://betterstack.com/community/guides/logging/logging-in-go/#error-logging-with-slog
	replace := func(_ []string, a slog.Attr) slog.Attr {
		if a.Key == "time" {
			return slog.Attr{}
		}
		return a
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: replace,
	})

	log := slog.New(handler)
	ctx := context.Background()

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Error(l.Error("error loading config", err))
	}

	conn := sqlite.NewDB(cfg.DSN)
	err := conn.Open()
	if err != nil {
		log.Error(l.Error("cannot open sqlite connection", err))
		panic(err)
	}
	log.Info(fmt.Sprintf("DSN is set to %v", cfg.DSN))

	serveApp := serve.NewServeProdApp(
		flashscoreUri,
		cfg.FlashscoreApiKey,
		cfg.DiscordUri,
		l.NewTelegramLogger(cfg.TelegramToken, cfg.TelegramAdmin, "serve"),
		serveImportUpcomingSchedule,
		serveImportFinishedSchedule,
	)

	// giftbox
	idGen := func() (string, error) {
		v7, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return v7.String(), nil
	}

	monitor := giftbox.NewTelegramMonitor(
		l.NewTelegramLogger(cfg.TelegramToken, cfg.TelegramAdmin, "giftbox"),
	)

	giftBoxServer, err := giftbox.NewServer(ctx, conn, idGen, cfg.GiftBoxApiKey, monitor)
	if err != nil {
		panic(fmt.Errorf("error starting gift box server: %v", err))
	}
	server := gracefulshutdown.NewServer(
		&http.Server{
			Addr:              ":8080",
			ReadHeaderTimeout: 10 * time.Second,
			Handler:           giftBoxServer,
		},
	)

	//pedroApp := telegram.NewPedro(
	//	conn,
	//	cfg.TelegramToken,
	//	cfg.AllowedUserIds,
	//	handler,
	//)
	//go pedroApp.Start()
	go serveApp.StartBackgroundJobs(ctx)
	if err := server.ListenAndServe(ctx); err != nil {
		log.Error(l.Error("didn't shut down gracefully", err))
	}
}
