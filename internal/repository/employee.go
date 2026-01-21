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

func (r *EmployeeRepository) GetAllEmployees(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee
	query := `SELECT e.id, e.nip, e.name, e.position, e.branch_office, e.employment_type, e.is_excluded, e.guaranteed_doorprize, a.present_at
			  FROM employees e
			  LEFT JOIN attendances a
			  	ON a.person_type = 'employee' AND a.person_id = e.id`
	err := r.db.SelectContext(ctx, &employees, query)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *EmployeeRepository) GetPresentEmployees(ctx context.Context) ([]model.Employee, error) {
	var employees []model.Employee
	query := `SELECT e.id, e.nip, e.name, e.position, e.branch_office, e.employment_type, e.is_excluded, e.guaranteed_doorprize, a.present_at
			  FROM employees e
			  JOIN attendances a
			  	ON a.person_type = 'employee' AND a.person_id = e.id
			  WHERE a.present_at IS NOT NULL`
	err := r.db.SelectContext(ctx, &employees, query)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *EmployeeRepository) GetEmployeeByName(ctx context.Context, name string) (*model.Employee, error) {
	var employee model.Employee
	query := `SELECT e.id, e.nip, e.name, e.position, e.branch_office, e.employment_type, e.is_excluded, e.guaranteed_doorprize, a.present_at
			  FROM employees e
			  LEFT JOIN attendances a
			  	ON a.person_type = 'employee' AND a.person_id = e.id
			  WHERE e.name = ?`
	err := r.db.GetContext(ctx, &employee, query, name)
	if err != nil {
		return nil, err
	}
	return &employee, nil
}

func (r *EmployeeRepository) UpdateEmployeePresentAt(ctx context.Context, empID int64, presentAt time.Time) error {
	query := `INSERT INTO attendances (person_type, person_id, present_at)
			  VALUES ('employee', ?, ?)`
	_, err := r.db.ExecContext(ctx, query, empID, presentAt)
	return err
}

func (r *EmployeeRepository) ResetAllAttendances(ctx context.Context) error {
	query := `DELETE FROM attendances WHERE person_type = 'employee'`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
