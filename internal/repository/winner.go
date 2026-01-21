package repository

import (
	"context"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/jmoiron/sqlx"
)

type WinnerRepository struct {
	db *sqlx.DB
}

func NewWinnerRepository(db *sqlx.DB) *WinnerRepository {
	return &WinnerRepository{
		db: db,
	}
}

func (r *WinnerRepository) AddWinners(ctx context.Context, winners []model.Winner) error {
	if len(winners) == 0 {
		return nil
	}

	query := `INSERT IGNORE INTO winners (
		employee_id,
		name,
		position,
		branch,
		employment_type,
		prize_type,
		round_id,
		round_label,
		won_at
	) VALUES (
		:employee_id,
		:name,
		:position,
		:branch,
		:employment_type,
		:prize_type,
		:round_id,
		:round_label,
		:won_at
	)`

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareNamedContext(ctx, query)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, winner := range winners {
		if _, err := stmt.ExecContext(ctx, winner); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *WinnerRepository) ListWinners(ctx context.Context) ([]model.Winner, error) {
	var winners []model.Winner
	query := `SELECT employee_id, name, position, branch, employment_type, prize_type, round_id, round_label, won_at
			  FROM winners
			  ORDER BY won_at ASC`
	err := r.db.SelectContext(ctx, &winners, query)
	if err != nil {
		return nil, err
	}
	return winners, nil
}

func (r *WinnerRepository) ListWinnersByType(ctx context.Context, prizeType string) ([]model.Winner, error) {
	var winners []model.Winner
	query := `SELECT employee_id, name, position, branch, employment_type, prize_type, round_id, round_label, won_at
			  FROM winners
			  WHERE prize_type = ?
			  ORDER BY won_at ASC`
	err := r.db.SelectContext(ctx, &winners, query, prizeType)
	if err != nil {
		return nil, err
	}
	return winners, nil
}

func (r *WinnerRepository) ResetWinners(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM winners")
	return err
}
