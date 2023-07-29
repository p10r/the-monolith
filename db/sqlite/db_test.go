package sqlite

import (
	"context"
	"pedro-go/domain/expect"
	"testing"
)

func TestDB(t *testing.T) {
	var ctx = context.TODO()
	t.Run("opens in-memory connection", func(t *testing.T) {
		checkConn(t, ctx, NewDB(":memory:"))
	})

	t.Run("opens file connection", func(t *testing.T) {
		checkConn(t, ctx, NewDB("../../local/local.db"))
	})

	t.Run("monitors errors", func(t *testing.T) {

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
