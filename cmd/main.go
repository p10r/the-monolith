package main

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"log"
	"pedro-go/telegram"
)

type PedroConfig struct {
	TelegramToken  string  `env:"TELEGRAM_TOKEN"`
	DSN            string  `env:"DSN"`
	AllowedUserIds []int64 `env:"ALLOWED_USER_IDS"`
}

func main() {
	ctx := context.Background()

	var cfg PedroConfig
	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatal(err)
	}

	log.Printf("DSN is set to %v", cfg.DSN)
	telegram.Pedro(cfg.TelegramToken, cfg.DSN, cfg.AllowedUserIds)
}
