package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"pedro-go/domain"
)

type SqliteEventMonitor struct {
	db *DB
}

type monitoringEvent struct {
	ID        int64
	EventType string
	Data      string
}

func NewEventMonitor(db *DB) SqliteEventMonitor {
	return SqliteEventMonitor{db: db}
}

func (m SqliteEventMonitor) Monitor(ctx context.Context, e domain.MonitoringEvent) {
	log.Printf("%v: %v", e.Name(), e)

	err := m.saveEvent(ctx, e)
	if err != nil {
		log.Printf("Could not save monitoring event %e to db. Reason: %v", e, err)
	}
}

func (m SqliteEventMonitor) saveEvent(ctx context.Context, e domain.MonitoringEvent) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func(tx *Tx) {
		if err := tx.Rollback(); err != nil {
			log.Printf("Could not rollback tx: %v", err)
		}
	}(tx)

	data, err := e.ToJSON()
	if err != nil {
		log.Printf("event_monitor: Could not serialize %e\n", e)
		return err
	}
	entity := monitoringEvent{
		ID:        0,
		EventType: e.Name(),
		Data:      string(data),
	}
	_, err = insertNewMonitoringEvent(ctx, tx, entity)
	if err != nil {
		return err
	}

	//TODO debug log here that event was saved

	_ = tx.Commit()
	return nil
}

func insertNewMonitoringEvent(ctx context.Context, tx *Tx, e monitoringEvent) (monitoringEvent, error) {
	result, err := tx.ExecContext(ctx, `
		INSERT INTO monitoring_events (event_type,data) VALUES (?,?)`,
		e.EventType, e.Data)
	if err != nil {
		return monitoringEvent{}, err
	}

	// Read back new e ID into caller argument.
	id, err := result.LastInsertId()
	if err != nil {
		return monitoringEvent{}, err
	}

	e.ID = id
	return e, nil
}

func (m SqliteEventMonitor) All(ctx context.Context) (domain.MonitoringEvents, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.MonitoringEvents{}, err
	}
	defer func(tx *Tx) {
		if err := tx.Rollback(); err != nil {
			log.Printf("Could not rollback tx: %v", err)
		}
	}(tx)

	entities, err := findMonitoringEvents(ctx, tx)
	if err != nil {
		return domain.MonitoringEvents{}, err
	}

	var artists []domain.MonitoringEvent
	for _, e := range entities {
		var err error
		switch e.EventType {

		case domain.ArtistFollowed{}.Name():
			var deserialized domain.ArtistFollowed
			err = json.Unmarshal([]byte(e.Data), &deserialized)
			artists = append(artists, deserialized)

		case domain.NewEventForArtist{}.Name():
			var deserialized domain.NewEventForArtist
			err = json.Unmarshal([]byte(e.Data), &deserialized)
			artists = append(artists, deserialized)
		}
		if err != nil {
			return domain.MonitoringEvents{}, fmt.Errorf("cannot deserialize monitoring event from db with ID %v", e.ID)
		}
	}

	return artists, err
}

func findMonitoringEvents(ctx context.Context, tx *Tx) ([]*monitoringEvent, error) {
	rows, err := tx.QueryContext(ctx, `
		SELECT id, event_type, data
		FROM monitoring_events 
		ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := rows.Close(); err != nil {
			log.Printf("Could not rollback tx: %v", err)
		}
	}(rows)

	entities := make([]*monitoringEvent, 0)
	for rows.Next() {
		var e monitoringEvent
		err := rows.Scan(
			&e.ID,
			&e.EventType,
			&e.Data,
		)
		if err != nil {
			return nil, err
		}
		entities = append(entities, &e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entities, nil
}
