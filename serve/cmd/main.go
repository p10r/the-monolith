package main

import (
	"context"
	"github.com/p10r/pedro/pkg/sqlite"
	"github.com/p10r/pedro/serve"
	"log"
	"os"
)

func main() {
	log.Println("Starting serve")

	conn := sqlite.NewDB(os.Getenv("DSN"))
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("DSN is set to %v", os.Getenv("DSN"))

	app := serve.NewServe(
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
	)

	_, err = app.Importer.ImportScheduledMatches(context.TODO())
	if err != nil {
		log.Fatal("Error when importing matches:", err)
	}
}
