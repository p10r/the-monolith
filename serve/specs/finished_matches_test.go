package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/p10r/monolith/serve/testutil"
	"testing"
)

func TestFinishedMatches(t *testing.T) {
	ctx := context.TODO()

	t.Run("sends scores to discord", func(t *testing.T) {
		f := newFixture(t, false, false)
		defer f.server.Close()

		err := f.importer.ImportFinishedMatches(ctx)
		assert.NoError(t, err)

		assert.NoError(t, err)

		discordReqs := *f.discordRequests
		msg := discordReqs[0]
		approvals.VerifyJSONBytes(t, testutil.PrettyPrinted(t, msg))
	})

	// run direnv allow . before running
	t.Run("run against prod", func(t *testing.T) {
		t.Skip()

		f := newFixture(t, true, false)
		defer f.server.Close()

		err := f.importer.ImportFinishedMatches(ctx)
		assert.NoError(t, err)

		// check discord
	})

}
