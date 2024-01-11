package main

import (
	"log"
	"os"
	"pedro-go/telegram"
)

var (
	token = os.Getenv("TELEGRAM_TOKEN")
	dsn   = os.Getenv("DSN")
)

func main() {
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is missing")
	}

	// GORM will create a new DB if the DSN doesn't match
	if dsn == "" {
		log.Fatal("DSN is missing")
	}

	telegram.Pedro(token, dsn)
}
