package model

import "time"

type Winner struct {
	EmployeeID     string    `json:"employee_id" db:"employee_id"`
	Name           string    `json:"name" db:"name"`
	Position       string    `json:"position" db:"position"`
	Branch         string    `json:"branch" db:"branch"`
	EmploymentType string    `json:"employment_type" db:"employment_type"`
	PrizeType      string    `json:"prize_type" db:"prize_type"`
	RoundID        string    `json:"round_id" db:"round_id"`
	RoundLabel     string    `json:"round_label" db:"round_label"`
	WonAt          time.Time `json:"won_at" db:"won_at"`
}
