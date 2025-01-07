package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestFinishedMatches(t *testing.T) {
	ctx := context.TODO()

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
