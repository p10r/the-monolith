package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	approvals "github.com/approvals/go-approval-tests"
	"github.com/approvals/go-approval-tests/reporters"
	"github.com/p10r/pedro/pedro/domain/expect"
	"github.com/p10r/pedro/serve/domain"
	"github.com/p10r/pedro/serve/testutil"
	"os"
	"sort"
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

	_, err := f.importer.ImportScheduledMatches(ctx)
	assert.NoError(t, err)

	t.Run("sends discord message", func(t *testing.T) {
		requests := *f.discordServer.Requests
		expect.Len(t, requests, 1)

		msg := newDiscordMessage(t, requests[0])
		approvals.VerifyJSONBytes(t, testutil.PrettyPrinted(t, msg))
	})
}

type matchWithoutID struct {
	HomeName  string
	AwayName  string
	StartTime int64
	Country   string
	League    string
}

func MatchesEqual(t *testing.T, got, want domain.Matches) {
	t.Helper()

	var gotten []matchWithoutID
	for _, match := range got {
		m := matchWithoutID{
			match.HomeName,
			match.AwayName,
			match.StartTime,
			match.Country,
			match.League,
		}
		gotten = append(gotten, m)
	}

	var wanted []matchWithoutID
	for _, match := range want {
		m := matchWithoutID{
			match.HomeName,
			match.AwayName,
			match.StartTime,
			match.Country,
			match.League,
		}
		wanted = append(wanted, m)
	}

	sort.Slice(gotten, func(i, j int) bool {
		return len(gotten[i].HomeName) > len(gotten[j].HomeName)
	})

	sort.Slice(wanted, func(i, j int) bool {
		return len(wanted[i].HomeName) > len(wanted[j].HomeName)
	})

	assert.Equal(t, gotten, wanted)
}

//TODO test what happens if two matches with the same timestamp are in db
//TODO show errors when DB is not there
