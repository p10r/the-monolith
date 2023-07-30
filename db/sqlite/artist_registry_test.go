package sqlite

import (
	"context"
	"pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestSqliteArtistRegistry(t *testing.T) {
	t.Run("adds and finds artists", func(t *testing.T) {
		var (
			ctx      = context.Background()
			recorder = domain.NewEventRecorder(nil)
			db       = MustOpenDB(t, &recorder)
			//db       = NewDB("../../local/local.db", &recorder)
			registry = NewArtistRegistry(db, &recorder)
		)
		db.Open()
		defer MustCloseDB(t, db)

		_, err := registry.Add(ctx, domain.NewArtist{Name: "Boys Noize"})
		expect.NoErr(t, err)

		got, err := registry.FindAll(ctx)
		want := domain.Artists{domain.Artist{Id: 1, Name: "Boys Noize"}}

		expect.NoErr(t, err)
		expect.SliceEqual(t, got, want)
	})

	//checks that artist doesn't already exist
}
