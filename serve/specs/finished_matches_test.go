package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/p10r/pedro/pedro/domain/expect"
	"github.com/p10r/pedro/serve/testutil"
	"testing"
)

func TestFinishedMatches(t *testing.T) {
	ctx := context.TODO()
	favs := []string{"Europe: Champions League Women - Play Offs"}
	f := newFixture(t, favs, false, false)
	defer f.flashscoreServer.Close()
	defer f.discordServer.Close()

	err := f.importer.ImportFinishedMatches(ctx)
	assert.NoError(t, err)

	t.Run("sends scores to discord", func(t *testing.T) {
		requests := *f.discordServer.Requests
		expect.Len(t, requests, 1)

		msg := newDiscordMessage(t, requests[0])
		approvals.VerifyJSONBytes(t, testutil.PrettyPrinted(t, msg))
	})

	t.Run("gets statistics from volleystation", func(t *testing.T) {

	})
}
