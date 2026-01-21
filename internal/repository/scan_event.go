package repository

import (
	"context"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/jmoiron/sqlx"
)

type ScanEventRepository struct {
	db *sqlx.DB
}

func NewScanEventRepository(db *sqlx.DB) *ScanEventRepository {
	return &ScanEventRepository{db: db}
}

func (r *ScanEventRepository) AddScanEvent(ctx context.Context, event model.ScanEvent) error {
	query := `INSERT INTO scan_events (scanner_id, name, scanned_at)
			  VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, event.ScannerID, event.Name, event.ScannedAt)
	return err
}

func (r *ScanEventRepository) ListRecentByScanner(ctx context.Context, scannerID int, limit int) ([]model.ScanEvent, error) {
	var events []model.ScanEvent
	query := `SELECT id, scanner_id, name, scanned_at
			  FROM scan_events
			  WHERE scanner_id = ?
			  ORDER BY scanned_at DESC, id DESC
			  LIMIT ?`
	err := r.db.SelectContext(ctx, &events, query, scannerID, limit)
	if err != nil {
		return nil, err
	}
	return events, nil
}
