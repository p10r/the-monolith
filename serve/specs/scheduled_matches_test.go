package specifications

import (
	"context"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/expect"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	r := approvals.UseReporter(reporters.NewIntelliJReporter())
	defer r.Close()

	approvals.UseFolder("testdata")
	os.Exit(m.Run())
}

func TestImportMatches(t *testing.T) {
	ctx := context.TODO()
	favs := []string{"Europe: Champions League - Play Offs", "USA: PVF Women"}

	f := newFixture(t, favs, false)
	defer f.flashscoreServer.Close()
	defer f.discordServer.Close()

	_, err := f.importer.ImportScheduledMatches(ctx)
	expect.NoErr(t, err)

	t.Run("imports today's matches to db", func(t *testing.T) {
		expected := domain.Matches{
			{
				HomeName:  "Trentino",
				AwayName:  "Jastrzebski",
				StartTime: 1714917600,
				Country:   "Europe",
				League:    "Champions League - Play Offs",
			},
			{
				HomeName:  "Resovia",
				AwayName:  "Zaksa",
				StartTime: 1714917600,
				Country:   "Europe",
				League:    "Champions League - Play Offs",
			},
			{
				HomeName:  "Grand Rapids Rise W",
				AwayName:  "San Diego Mojo W",
				StartTime: 1714939200,
				Country:   "USA",
				League:    "PVF Women",
			},
		}
		expect.MatchStoreContains(t, f.store, expected)
	})

	t.Run("sends discord message", func(t *testing.T) {
		requests := *f.discordServer.Requests
		expect.Len(t, requests, 1)

		msg := newDiscordMessage(t, requests[0])
		approvals.VerifyJSONBytes(t, prettyPrinted(t, msg))
	})

	// Make sure to:
	// 1. remove t.Skip()
	// 2. direnv allow . && go test specs/scheduled_matches_test.go
	t.Run("run against real discord", func(t *testing.T) {
		t.Skip()

		_, _ = newFixture(t, favs, true).importer.ImportScheduledMatches(ctx)
	})

	//TODO sends to discord even if db is not available
}

//TODO test what happens if two matches with the same timestamp are in db
//TODO show errors when DB is not there
