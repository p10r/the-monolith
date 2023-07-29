package sqlite

import (
	"context"
	"errors"
	"pedro-go/domain"
	"pedro-go/domain/expect"
	"testing"
)

func TestDB(t *testing.T) {
	var (
		ctx      = context.TODO()
		events   = domain.Events{}
		recorder = domain.NewEventRecorder(events)
	)
	t.Run("opens in-memory connection", func(t *testing.T) {
		checkConn(t, ctx, NewDB(":memory:", &recorder))
	})

	t.Run("opens file connection", func(t *testing.T) {
		checkConn(t, ctx, NewDB("../../local/local.db", &recorder))
	})

	t.Run("monitors errors", func(t *testing.T) {
		db := NewDB("", &recorder)

		err := db.Open()
		expect.Err(t, err)

		want := domain.Events{domain.ErrEvent{Err: errors.New("dsn required").Error()}}
		expect.SliceEqual(t, recorder.Events, want)
	})
}

func checkConn(t *testing.T, ctx context.Context, db *DB) {
	err := db.Open()
	expect.NoErr(t, err)

	tx, err := db.BeginTx(ctx, nil)
	expect.NoErr(t, err)

	rows, err := tx.Query("SELECT 1;")
	expect.NoErr(t, err)

	var queryResult int
	for rows.Next() {
		err = rows.Scan(&queryResult)
		expect.NoErr(t, err)
	}

	expect.Equal(t, queryResult, 1)
}

func MustOpenDB(tb testing.TB, recorder domain.EventRecorder) *DB {
	tb.Helper()

	// Write to an in-memory database by default.
	// If the -dump flag is set, generate a temp file for the database.
	dsn := ":memory:"
	db := NewDB(dsn, recorder)
	if err := db.Open(); err != nil {
		tb.Fatal(err)
	}
	return db
}
