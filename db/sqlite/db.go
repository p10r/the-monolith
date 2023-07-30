// Package sqlite is taken from Ben Johnson's wtf
// https://github.com/benbjohnson/wtf/blob/030fcb0d5ff21b56fba09564fbba7e691ae50886/sqlite/sqlite.go
package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"pedro-go/domain"
	"sort"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection.
type DB struct {
	db     *sql.DB
	ctx    context.Context // background context
	cancel func()          // cancel background context

	// Datasource name.
	DSN string

	// Returns the current time. Defaults to time.Now().
	// Can be mocked for tests.
	Now func() time.Time

	// tracks all events and takes care of logging and storing them
	Events domain.EventRecorder
}

func NewDB(dsn string, recorder domain.EventRecorder) *DB {
	db := &DB{
		DSN:    dsn,
		Now:    time.Now,
		Events: recorder,
	}
	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) Open() (err error) {
	// Ensure a DSN is set before attempting to open the database.
	if db.DSN == "" {
		err := "dsn required"
		db.Events.Record(domain.ErrEvent{Err: err})
		return fmt.Errorf(err)
	}
	db.Events.Record(domain.DbEvent{Msg: fmt.Sprintf("Setting up datasource '%v'", db.DSN)})

	// Make the parent directory unless using an in-memory db.
	if db.DSN != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(db.DSN), 0700); err != nil {
			db.Events.Record(domain.ErrEvent{Err: err.Error()})
			return err
		}
	}

	// Connect to the database.
	if db.db, err = sql.Open("sqlite3", db.DSN); err != nil {
		return err
	}

	db.Events.Record(domain.DbEvent{Msg: "Setting journal_mode = wal"})
	// WAL allows multiple readers to operate while data is being written.
	if _, err := db.db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		db.Events.Record(domain.ErrEvent{Err: err.Error()})
		return fmt.Errorf("enable wal: %w", err)
	}

	db.Events.Record(domain.DbEvent{Msg: "Setting foreign_keys = ON"})
	if _, err := db.db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		db.Events.Record(domain.ErrEvent{Err: err.Error()})
		return fmt.Errorf("foreign keys pragma: %w", err)
	}

	db.Events.Record(domain.DbEvent{Msg: "Running migrations"})
	if err := db.migrate(); err != nil {
		db.Events.Record(domain.ErrEvent{Err: err.Error()})
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

// Close closes the database connection.
func (db *DB) Close() error {
	// Cancel background context.
	db.cancel()

	// Close database.
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

// BeginTx starts a transaction and returns a wrapper Tx type. This type
// provides a reference to the database and a fixed timestamp at the start of
// the transaction. The timestamp allows us to mock time during tests as well.
func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		db.Events.Record(domain.ErrEvent{Err: err.Error()})
		return nil, err
	}

	// Return wrapper Tx that includes the transaction start time.
	return &Tx{
		Tx:  tx,
		db:  db,
		now: db.Now().UTC().Truncate(time.Second),
	}, nil
}

// Tx wraps the SQL Tx object to provide a timestamp at the start of the transaction.
type Tx struct {
	*sql.Tx
	db  *DB
	now time.Time
}

// NullTime represents a helper wrapper for time.Time. It automatically converts
// time fields to/from RFC 3339 format. Also supports NULL for zero time.
type NullTime time.Time

// Scan reads a time value from the database.
func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(n) = time.Time{}
		return nil
	} else if value, ok := value.(string); ok {
		*(*time.Time)(n), _ = time.Parse(time.RFC3339, value)
		return nil
	}

	return fmt.Errorf("NullTime: cannot scan to time.Time: %T", value)
}

// Value formats a time value for the database.
func (n *NullTime) Value() (driver.Value, error) {
	if n == nil || (*time.Time)(n).IsZero() {
		return nil, nil
	}
	return (*time.Time)(n).UTC().Format(time.RFC3339), nil
}

//go:embed migrations/*.sql
var migrationFS embed.FS

// migrate sets up migration tracking and executes pending migration files.
//
// Migration files are embedded in the sqlite/migration folder and are executed
// in lexigraphical order.
//
// Once a migration is run, its name is stored in the 'migrations' table so it
// is not re-executed. Migrations run in a transaction to prevent partial
// migrations.
func (db *DB) migrate() error {
	// Ensure the 'migrations' table exists, so we don't duplicate migrations.
	if _, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	// Read migration files from our embedded file system.
	// This uses Go 1.16's 'embed' package.
	names, err := fs.Glob(migrationFS, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	// Loop over all migration files and execute them in order.
	for _, name := range names {
		err := db.migrateFile(name)
		if err != nil {
			err := fmt.Errorf("migration error: name=%q err=%w", name, err)
			return err
		}
	}
	return nil
}

// migrate runs a single migration file within a transaction. On success, the
// migration file name is saved to the "migrations" table to prevent re-running.
func (db *DB) migrateFile(name string) error {

	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Ensure migration has not already been run.
	var n int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = ?`, name).Scan(&n); err != nil {
		return err
	} else if n != 0 {
		return nil // already run migration, skip
	}

	// Read and execute migration file.
	if buf, err := fs.ReadFile(migrationFS, name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	//var a int
	//tx.Exec(`INSERT INTO artists (name) VALUES ('tests')`)
	//tx.QueryRow(`SELECT COUNT(*) FROM artists`).Scan(&a)
	//fmt.Printf("a is %v \n", a)

	// Insert record into migrations to prevent re-running migration.
	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES (?)`, name); err != nil {
		return err
	}

	return tx.Commit()
}
