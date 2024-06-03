package specifications

import (
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/serve/db"
	"github.com/p10r/pedro/serve/discord"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/flashscore"
	"github.com/p10r/pedro/serve/testutil"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type fixture struct {
	flashscoreServer *httptest.Server
	discordServer    *testutil.DiscordServer
	importer         *domain.MatchImporter
	store            *db.MatchStore
}

func newFixture(
	t *testing.T,
	favLeagues []string,
	runAgainstDiscord bool,
	runAgainstFlashscore bool,
) fixture {
	log := l.NewTextLogger().With(slog.String("app", "serve"))

	apiKey := "random_api_key"

	var flashscoreServer *httptest.Server
	var fsClient *flashscore.Client
	if runAgainstFlashscore {
		key := os.Getenv("FLASHSCORE_API_KEY")
		if key == "" {
			t.Fatalf("No FLASHSCORE_API_KEY set. Run direnv allow .")
		}

		fsClient = flashscore.NewClient("https://flashscore.p.rapidapi.com", key, log)
	} else {
		flashscoreServer = testutil.NewFlashscoreServer(t, apiKey)
		fsClient = flashscore.NewClient(flashscoreServer.URL, apiKey, log)
	}

	var discordServer *testutil.DiscordServer
	var discordClient *discord.Client
	if runAgainstDiscord {
		uri := os.Getenv("DISCORD_URI")
		if uri == "" {
			t.Fatalf("No DISCORD_URI set. Run direnv allow .")
		}
		discordClient = discord.NewClient(uri, log)
	} else {
		discordServer = testutil.NewDiscordServer(t, log)
		discordClient = discord.NewClient(discordServer.URL, log)
	}

	may28th := func() time.Time {
		return time.Date(2024, 5, 28, 0, 0, 0, 0, time.UTC)
	}

	matchStore := db.NewMatchStore(testutil.MustOpenDB(t))
	importer := domain.NewMatchImporter(
		matchStore,
		fsClient,
		discordClient,
		favLeagues,
		may28th,
		log,
	)
	return fixture{
		flashscoreServer,
		discordServer,
		importer,
		matchStore,
	}
}
