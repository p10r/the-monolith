package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
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

	f := newFixture(t, false, false)
	defer f.server.Close()

	_, err := f.importer.ImportScheduledMatches(ctx)
	assert.NoError(t, err)
}
