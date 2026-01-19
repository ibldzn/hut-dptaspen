package repository

import (
	"context"
	"time"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/jmoiron/sqlx"
)

type EmployeeRepository struct {
	db *sqlx.DB
}

func NewEmployeeRepository(db *sqlx.DB) *EmployeeRepository {
	return &EmployeeRepository{
		db: db,
	}
}

func (r *EmployeeRepository) GetPresentEmployees(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee
	query := `SELECT id, name, position, branch_office, employment_type, is_excluded, guaranteed_doorprize, present_at
			  FROM employees
			  WHERE present_at IS NULL`
	err := r.db.SelectContext(ctx, &employees, query)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *EmployeeRepository) GetEmployeeByName(ctx context.Context, name string) (*model.Employee, error) {
	var employee model.Employee
	query := `SELECT id, name, position, branch_office, employment_type, is_excluded, guaranteed_doorprize, present_at
			  FROM employees
			  WHERE name = ?`
	err := r.db.GetContext(ctx, &employee, query, name)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) UpdateEmployeePresentAt(ctx context.Context, empID int64, presentAt time.Time) error {
	query := `UPDATE employees SET present_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, presentAt, empID)
	return err
}
