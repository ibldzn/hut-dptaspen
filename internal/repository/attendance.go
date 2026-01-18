package repository

import "github.com/jmoiron/sqlx"

type AttendanceRepository struct {
	db *sqlx.DB
}

func NewAttendanceRepository(db *sqlx.DB) *AttendanceRepository {
	return &AttendanceRepository{
		db: db,
	}
}
