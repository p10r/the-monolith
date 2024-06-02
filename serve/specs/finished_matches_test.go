package specifications

import (
	"context"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/p10r/pedro/serve/expect"
	"testing"
)

func TestFinishedMatches(t *testing.T) {
	ctx := context.TODO()
	favs := []string{"Europe: Champions League Women - Play Offs"}
	f := newFixture(t, favs, false)
	defer f.flashscoreServer.Close()
	defer f.discordServer.Close()

	err := f.importer.ImportFinishedMatches(ctx)
	expect.NoErr(t, err)

	t.Run("sends scores to discord", func(t *testing.T) {
		requests := *f.discordServer.Requests
		expect.Len(t, requests, 1)

		msg := newDiscordMessage(t, requests[0])
		approvals.VerifyJSONBytes(t, prettyPrinted(t, msg))
	})

	t.Run("gets statistics from volleystation", func(t *testing.T) {

	})

	// Make sure to:
	// 1. remove t.Skip()
	// 2. direnv allow . && go test serve/specs/finished_matches_test.go
	t.Run("run against real discord", func(t *testing.T) {
		t.Skip()

		_ = newFixture(t, favs, true).importer.ImportFinishedMatches(ctx)
	})
}
