package model

import "time"

type ScanEvent struct {
	ID        int64     `json:"id" db:"id"`
	ScannerID int       `json:"scanner_id" db:"scanner_id"`
	Name      string    `json:"name" db:"name"`
	ScannedAt time.Time `json:"scanned_at" db:"scanned_at"`
}
