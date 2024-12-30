package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/p10r/pedro/pedro/domain/expect"
	"github.com/p10r/pedro/serve/testutil"
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

	f := newFixture(t, favs, false, false)
	defer f.flashscoreServer.Close()
	defer f.discordServer.Close()
	defer f.superLegaWebsite.Close()

	_, err := f.importer.ImportScheduledMatches(ctx)
	assert.NoError(t, err)

	t.Run("sends discord message", func(t *testing.T) {
		requests := *f.discordServer.Requests
		expect.Len(t, requests, 1)

		msg := newDiscordMessage(t, requests[0])
		approvals.VerifyJSONBytes(t, testutil.PrettyPrinted(t, msg))
	})
}
