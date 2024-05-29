package main

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"log"
	"pedro-go/pedro/telegram"
)

type Config struct {
	TelegramToken  string  `env:"TELEGRAM_TOKEN"`
	DSN            string  `env:"DSN"`
	AllowedUserIds []int64 `env:"ALLOWED_USER_IDS"`
}

func main() {
	ctx := context.Background()

	var cfg Config
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatal(err)
	}

	telegram.Pedro(cfg.TelegramToken, cfg.DSN, cfg.AllowedUserIds).Start()
}
