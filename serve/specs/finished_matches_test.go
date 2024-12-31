package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestFinishedMatches(t *testing.T) {
	ctx := context.TODO()
	favs := []string{
		"Europe: Champions League Women - Play Offs",
		"Poland: PlusLiga",
		"Italy: SuperLega",
	}

	// run direnv allow . before running
	t.Run("run against prod", func(t *testing.T) {
		t.Skip()

		f := newFixture(t, favs, true, false)
		defer f.server.Close()

		err := f.importer.ImportFinishedMatches(ctx)
		assert.NoError(t, err)

		// check discord
	})

}
