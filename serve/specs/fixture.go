package specifications

import (
	"github.com/p10r/pedro/pkg/l"
	"github.com/p10r/pedro/serve/discord"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/flashscore"
	"github.com/p10r/pedro/serve/statistics"
	"github.com/p10r/pedro/serve/testutil"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type fixture struct {
	server   *httptest.Server
	importer *domain.MatchImporter
}

func newFixture(
	t *testing.T,
	favLeagues []string,
	runAgainstDiscord bool,
	runAgainstFlashscore bool,
) fixture {
	log := l.NewTextLogger().With(slog.String("app", "serve"))

	apiKey := "random_api_key"

	mux := http.NewServeMux()
	mux.Handle("GET /flashscore/v1/events/list", testutil.NewFlashscoreServer(t, apiKey))
	mux.Handle("POST /discord", testutil.NewDiscordServer(t, log))
	mux.Handle("GET /plusliga", testutil.NewPlusLigaServer(
		t, testutil.MustReadFile(t, "../testdata/statistics/plusliga.html"),
	))
	mux.Handle("GET /superlega", testutil.NewSuperLegaServer(
		t, testutil.MustReadFile(t, "../testdata/statistics/superlega-italy-m.html"),
	))
	server := httptest.NewServer(mux)

	var fsClient *flashscore.Client
	if runAgainstFlashscore {
		key := os.Getenv("FLASHSCORE_API_KEY")
		if key == "" {
			t.Fatalf("No FLASHSCORE_API_KEY set. Run direnv allow .")
		}

		fsClient = flashscore.NewClient("https://flashscore.p.rapidapi.com", key, log)
	} else {
		fsClient = flashscore.NewClient(server.URL+"/flashscore", apiKey, log)
	}

	var discordClient *discord.Client
	if runAgainstDiscord {
		uri := os.Getenv("DISCORD_URI")
		if uri == "" {
			t.Fatalf("No DISCORD_URI set. Run direnv allow .")
		}
		discordClient = discord.NewClient(uri, log)
	} else {
		discordClient = discord.NewClient(server.URL+"/discord", log)
	}

	aggr := statistics.NewAggregator(server.URL+"/plusliga", server.URL+"/superlega", log)

	may28th := func() time.Time {
		return time.Date(2024, 5, 28, 0, 0, 0, 0, time.UTC)
	}

	return fixture{
		server,
		domain.NewMatchImporter(
			fsClient,
			discordClient,
			aggr,
			favLeagues,
			may28th,
			log,
		),
	}
}
