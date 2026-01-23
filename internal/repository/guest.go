package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/jmoiron/sqlx"
)

type GuestRepository struct {
	db *sqlx.DB
}

func NewGuestRepository(db *sqlx.DB) *GuestRepository {
	return &GuestRepository{
		db: db,
	}
}

func (r *GuestRepository) ListGuests(ctx context.Context) ([]model.Guest, error) {
	var guests []model.Guest
	query := "SELECT id, name, `table`, present_at FROM guests WHERE present_at IS NOT NULL ORDER BY present_at ASC"
	err := r.db.SelectContext(ctx, &guests, query)
	if err != nil {
		return nil, err
	}
	return guests, nil
}

func (r *GuestRepository) GetGuestByName(ctx context.Context, name string) (*model.Guest, error) {
	var guest model.Guest
	query := "SELECT id, name, `table`, present_at FROM guests WHERE name = ?"
	err := r.db.GetContext(ctx, &guest, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &guest, nil
}

func (r *GuestRepository) AddGuest(ctx context.Context, name string, presentAt time.Time) error {
	query := `INSERT INTO guests (name, present_at)
			  VALUES (?, ?)`
	_, err := r.db.ExecContext(ctx, query, name, presentAt)
	return err
}

func (r *GuestRepository) ResetGuests(ctx context.Context) error {
	query := `DELETE FROM guests`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
