package giftbox

import (
	"context"
	"github.com/p10r/pedro/pkg/sqlite"
	"log"
)

type GiftRepository struct {
	db *sqlite.DB
}

func NewGiftRepository(conn *sqlite.DB) *GiftRepository {
	return &GiftRepository{conn}
}

type dbEntity struct {
	// This will never grow big, so that's fine
	id       int
	lookupId string
	giftType string
	// Either 0 or 1. Sqlite doesn't have a bool type
	redeemed int
	json     string
}

func (r *GiftRepository) Save(ctx context.Context, gift Gift) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	//nolint:errcheck
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
	INSERT INTO gifts (
		lookup_id,
		gift_type,
		redeemed,
		json
	) VALUES (?,?,?,?)`, gift.ID, gift.Type, gift.Redeemed, "{}")
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	log.Printf("Inserted new record %v for %v", id, gift)
	_ = tx.Commit()
	return nil
}

func (r *GiftRepository) All(ctx context.Context) (Gifts, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Gifts{}, err
	}
	//nolint:errcheck
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, `
		SELECT 
			id,
			lookup_id,
			gift_type,
			redeemed,
			json
		FROM gifts 
		ORDER BY id`)
	if err != nil {
		return Gifts{}, err
	}
	defer rows.Close()

	var gifts Gifts
	for rows.Next() {
		var e dbEntity
		err := rows.Scan(
			&e.id,
			&e.lookupId,
			&e.giftType,
			&e.redeemed,
			&e.json,
		)
		if err != nil {
			return Gifts{}, err
		}

		gifts = append(gifts, e.toGift())
	}
	if err := rows.Err(); err != nil {
		return Gifts{}, err
	}
	_ = tx.Commit()
	return gifts, err
}

func (e dbEntity) toGift() Gift {
	var redeemed bool
	if e.redeemed == 0 {
		redeemed = false
	} else {
		redeemed = true
	}

	var giftType GiftType
	switch e.giftType {
	case string(TypeSweet):
		giftType = TypeSweet
	case string(TypeWish):
		giftType = TypeWish
	}

	g := Gift{
		ID:       e.lookupId,
		Type:     giftType,
		Redeemed: redeemed,
	}
	return g
}

func (r *GiftRepository) SetRedeemedFlag(
	ctx context.Context,
	lookupId string,
	redeemed bool,
) (int64, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	//nolint:errcheck
	defer tx.Rollback()

	var dbBool int
	if redeemed {
		dbBool = 1
	} else {
		dbBool = 0
	}

	result, err := tx.ExecContext(ctx, `
		UPDATE gifts
		SET redeemed = ?
		WHERE lookup_id = ?`,
		dbBool, lookupId)
	if err != nil {
		return 0, err
	}

	_ = tx.Commit()

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
