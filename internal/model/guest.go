package model

import "time"

type Guest struct {
	ID        int64      `json:"ID" db:"id"`
	NamaTamu  string     `json:"NAMA_TAMU" db:"name"`
	Meja      *string    `json:"MEJA" db:"table"`
	PresentAt *time.Time `json:"PRESENT_AT" db:"present_at"`
}
