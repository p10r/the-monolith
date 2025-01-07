package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	"log"
	"testing"
)

func TestProd(t *testing.T) {
	ctx := context.TODO()

	// Make sure to:
	// 1. remove t.Skip()
	// 2. direnv allow . && go test specs/scheduled_matches_test.go
	t.Run("run against real discord", func(t *testing.T) {
		t.Skip()

		f := newFixture(t, true, false)

		_, _ = f.importer.ImportScheduledMatches(ctx)
	})

	// Make sure to:
	// 1. remove t.Skip()
	// 2. direnv allow . && go test specs/scheduled_matches_test.go
	t.Run("fetch real schedule", func(t *testing.T) {
		t.Skip()
		f := newFixture(t, false, true)

		matches, err := f.importer.ImportScheduledMatches(ctx)
		assert.NoError(t, err)
		for _, match := range matches {
			log.Println(match)
		}
	})

}
