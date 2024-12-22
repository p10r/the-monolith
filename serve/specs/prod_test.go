package specifications

import (
	"context"
	"github.com/alecthomas/assert/v2"
	"log"
	"testing"
)

func TestProd(t *testing.T) {
	ctx := context.TODO()
	favs := []string{
		"Italy: SuperLega",
		"Italy: SuperLega - Play Offs",
		"Italy: Coppa Italia A1",
		"Italy: Coppa Italia A1 Women",
		"Italy: Serie A1 Women",
		"Italy: Serie A1 Women - Playoffs",
		"Poland: PlusLiga",
		"Poland: PlusLiga - Play Offs",
		"France: Ligue A - Play Offs",
		"France: Ligue A",
		"Russia: Super League - Play Offs",
		"Russia: Super League",
		"Russia: Russia Cup",
		"World: Nations League",
		"World: Nations League - Play Offs",
		"World: Nations League Women",
		"World: Nations League Women - Play Offs",
		"World: Pan-American Cup",
		"World: World Championship - First round",
		"World: World Championship - Second round",
		"World: World Championship - Play Offs",
		"World: World Championship Women - First round",
		"Germany: VBL Supercup",
		"Germany: 1. Bundesliga",
		"Germany: 1. Bundesliga - Losers stage",
		"Germany: 1. Bundesliga - Winners stage",
		"Germany: 1. Bundesliga - Play Offs",
		"Germany: DVV Cup",
		"Turkey: Sultanlar Ligi Women",
		"Turkey: Sultanlar Ligi Women - Play Offs",
		"Turkey: Efeler Ligi",
		"TURKEY: Efeler Ligi - Play Offs",
		"Turkey: Efeler Ligi - 5th-8th places",
		"Europe: Champions League",
		"Europe: Champions League Women",
		"Europe: Champions League Women - Play Offs",
		"Europe: Champions League - Play Offs",
		"Europe: CEV Cup",
		"Europe: European Championships Women",
		"Europe: European Championships",
	}

	// Make sure to:
	// 1. remove t.Skip()
	// 2. direnv allow . && go test specs/scheduled_matches_test.go
	t.Run("run against real discord", func(t *testing.T) {
		t.Skip()

		f := newFixture(t, favs, true, false)

		_, _ = f.importer.ImportScheduledMatches(ctx)
	})

	// Make sure to:
	// 1. remove t.Skip()
	// 2. direnv allow . && go test specs/scheduled_matches_test.go
	t.Run("fetch real schedule", func(t *testing.T) {
		t.Skip()
		f := newFixture(t, favs, false, true)

		matches, err := f.importer.ImportScheduledMatches(ctx)
		assert.NoError(t, err)
		for _, match := range matches {
			log.Println(match)
		}
	})

}
