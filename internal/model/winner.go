package model

import "time"

type Winner struct {
	EmployeeID     string    `json:"employee_id"`
	Name           string    `json:"name"`
	Position       string    `json:"position"`
	Branch         string    `json:"branch"`
	EmploymentType string    `json:"employment_type"`
	PrizeType      string    `json:"prize_type"`
	RoundID        string    `json:"round_id"`
	RoundLabel     string    `json:"round_label"`
	WonAt          time.Time `json:"won_at"`
}
