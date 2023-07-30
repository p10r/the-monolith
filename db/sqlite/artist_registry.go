package sqlite

import (
	"context"
	_ "github.com/mattn/go-sqlite3"
	"pedro-go/domain"
)

type ArtistRegistry struct {
	db            *DB
	EventRecorder domain.EventRecorder
}

func NewArtistRegistry(db *DB, recorder domain.EventRecorder) *ArtistRegistry {
	return &ArtistRegistry{db: db, EventRecorder: recorder}
}

func (r ArtistRegistry) FindAll(ctx context.Context) (domain.Artists, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.EventRecorder.Record(domain.DbEvent{err.Error()})
		return nil, err
	}
	defer tx.Rollback()

	all, err := findAll(ctx, tx)
	if err != nil {
		r.EventRecorder.Record(domain.DbEvent{err.Error()})
		return nil, err
	}

	return all, nil
}

func findAll(ctx context.Context, tx *Tx) (domain.Artists, error) {
	query := "SELECT * FROM artists"

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	artists := make(domain.Artists, 0)
	for rows.Next() {
		var artist domain.Artist

		if err := rows.Scan(
			&artist.Id,
			&artist.Name,
		); err != nil {
			return nil, err
		}

		artists = append(artists, artist)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return artists, nil
}

func (r ArtistRegistry) Add(ctx context.Context, artist domain.NewArtist) (domain.Artist, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		r.EventRecorder.Record(domain.DbEvent{err.Error()})
		return domain.Artist{}, err
	}

	//_, err = tx.ExecContext(ctx, `CREATE TABLE artists(id INTEGER PRIMARY KEY AUTOINCREMENT)`, artist.Name)
	//if err != nil {
	//	r.EventRecorder.Record(DbEvent{Msg: err.Error()})
	//	//return domain.Artist{}, err
	//}

	result, err := tx.ExecContext(ctx, `INSERT INTO artists (name) VALUES (?)`, artist.Name)
	if err != nil {
		r.EventRecorder.Record(domain.DbEvent{Msg: err.Error()})
		return domain.Artist{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.EventRecorder.Record(domain.DbEvent{Msg: err.Error()})
		return domain.Artist{}, err
	}

	err = tx.Commit()
	if err != nil {
		r.EventRecorder.Record(domain.DbEvent{Msg: err.Error()})
		return domain.Artist{}, err
	}

	return domain.Artist{Id: domain.Id(id), Name: artist.Name}, nil
}
