package serve

import (
	"context"
	"github.com/p10r/pedro/pkg/sqlite"
	"github.com/p10r/pedro/serve/db"
	"github.com/p10r/pedro/serve/discord"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/flashscore"
	"log"
	"time"
)

type Serve struct {
	Importer *domain.MatchImporter
}

// NewServe wires Serve App together.
// Expects an already opened connection.
func NewServe(
	conn *sqlite.DB,
	flashscoreUri, flashscoreApiKey, discordUri string,
	favouriteLeagues []string,
) Serve {
	if flashscoreUri == "" {
		log.Fatal("flashscoreUri has not been set")
	}
	if flashscoreApiKey == "" {
		log.Fatal("flashscoreApiKey has not been set")
	}
	if discordUri == "" {
		log.Fatal("DISCORD_URI has not been set")
	}

	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewMatchStore(conn)
	flashscoreClient := flashscore.NewClient(flashscoreUri, flashscoreApiKey)
	discordClient := discord.NewClient(discordUri)

	now := func() time.Time { return time.Now() }

	importer := domain.NewMatchImporter(
		store,
		flashscoreClient,
		discordClient,
		favouriteLeagues,
		now,
	)

	return Serve{importer}
}

func (s Serve) Start(_ context.Context) {
	// TODO
}
