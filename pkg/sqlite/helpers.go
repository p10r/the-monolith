package sqlite

import "testing"

// MustOpenDB returns a new, open DB. Fatal on error.
func MustOpenDB(tb testing.TB) *DB {
	tb.Helper()

	// Write to an in-memory database by default.
	// If the -dump flag is set, generate a temp file for the database.
	dsn := ":memory:"

	instance := NewDB(dsn)
	if err := instance.Open(); err != nil {
		tb.Fatal(err)
	}
	return instance
}

// MustCloseDB closes the DB. Fatal on error.
func MustCloseDB(tb testing.TB, db *DB) {
	tb.Helper()
	if err := db.Close(); err != nil {
		tb.Fatal(err)
	}
}
