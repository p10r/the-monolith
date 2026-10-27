package specifications

import (
	"cmp"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/p10r/monolith/pkg/l"
	"github.com/p10r/monolith/serve/discord"
	"github.com/p10r/monolith/serve/domain"
	"github.com/p10r/monolith/serve/flashscore"
	"github.com/p10r/monolith/serve/testutil"
)

type fixture struct {
	server          *httptest.Server
	importer        *domain.MatchImporter
	discordRequests *[]discord.Message
}

func newFixture(
	t *testing.T,
	runAgainstDiscord bool,
	runAgainstFlashscore bool,
) fixture {
	log := l.NewTextLogger().With(slog.String("app", "serve"))

	apiKey := "random_api_key"

	discordReqs := []discord.Message{}

	mux := http.NewServeMux()
	mux.Handle("GET /flashscore/v1/events/list", testutil.NewFlashscoreServer(t, apiKey))
	discordReqRecorder := func(r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("Cannot read body")
			return
		}
		var discordMsg discord.Message
		if err := json.Unmarshal(body, &discordMsg); err != nil {
			t.Fatalf("Cannot parse discord message")
		}

		// Sort, to always have the same order in the message to help approval tests
		slices.SortFunc(discordMsg.Embeds[0].Fields, func(a, b discord.Fields) int {
			return cmp.Compare(a.Name, b.Name)
		})

		discordReqs = append(discordReqs, discordMsg)
		defer r.Body.Close()
	}
	mux.Handle("POST /discord", testutil.NewDiscordServer(t, log, discordReqRecorder))
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

	may28th := func() time.Time {
		return time.Date(2024, 5, 28, 0, 0, 0, 0, time.UTC)
	}

	return fixture{
		server,
		domain.NewMatchImporter(
			fsClient,
			discordClient,
			may28th,
			log,
		),
		&discordReqs,
	}
}
