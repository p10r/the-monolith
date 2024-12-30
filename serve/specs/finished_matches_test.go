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
	favs := []string{
		"Europe: Champions League Women - Play Offs",
		"Poland: PlusLiga",
		"Italy: SuperLega",
	}

	t.Run("sends scores to discord", func(t *testing.T) {
		f := newFixture(t, favs, false, false)
		defer f.flashscoreServer.Close()
		defer f.discordServer.Close()
		defer f.plusLigaWebsite.Close()
		defer f.superLegaWebsite.Close()

		err := f.importer.ImportFinishedMatches(ctx)
		assert.NoError(t, err)

		requests := *f.discordServer.Requests
		expect.Len(t, requests, 1)

		msg := newDiscordMessage(t, requests[0])
		approvals.VerifyJSONBytes(t, testutil.PrettyPrinted(t, msg))
	})

	// run direnv allow . before running
	t.Run("run against prod", func(t *testing.T) {
		t.Skip()

		f := newFixture(t, favs, true, false)
		defer f.flashscoreServer.Close()
		defer f.plusLigaWebsite.Close()
		defer f.superLegaWebsite.Close()

		err := f.importer.ImportFinishedMatches(ctx)
		assert.NoError(t, err)

		// check discord
	})

}
