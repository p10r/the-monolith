package main

import (
	"context"
	"fmt"
	"github.com/p10r/pedro/pkg/sqlite"
	"github.com/p10r/pedro/serve"
	"log/slog"
	"os"
)

func main() {
	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	log := slog.New(logHandler).With(slog.String("app", "serve"))

	conn := sqlite.NewDB(os.Getenv("DSN"))
	err := conn.Open()
	if err != nil {
		log.Error("cannot open db conn", slog.Any("error", err))
		panic(err)
	}

	log.Info(fmt.Sprintf("DSN is set to %v", os.Getenv("DSN")))

	app := serve.NewServeApp(
		conn,
		"https://flashscore.p.rapidapi.com",
		os.Getenv("FLASHSCORE_API_KEY"),
		os.Getenv("DISCORD_URI"),
		[]string{
			"Italy: SuperLega",
			"Italy: SuperLega - Play Offs",
			"World: Nations League",
			"World: Nations League - Play Offs",
		},
		logHandler,
	)

	_, err = app.Importer.ImportScheduledMatches(context.TODO())
	if err != nil {
		log.Error("Error when importing matches:", slog.Any("error", err))
		panic(err)
	}
}
