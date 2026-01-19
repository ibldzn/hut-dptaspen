package services

import (
	"context"
	"time"

	"github.com/ibldzn/spinner-hut/internal/model"
	"github.com/ibldzn/spinner-hut/internal/repository"
)

type EmployeeService struct {
	EmpRepository *repository.EmployeeRepository
}

func NewEmployeeService(empRepo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		EmpRepository: empRepo,
	}
}

func (s *EmployeeService) GetPresentEmployees(ctx context.Context) ([]model.Employee, error) {
	return s.EmpRepository.GetPresentEmployees(ctx)
}

func (s *EmployeeService) GetEmployeeByName(ctx context.Context, name string) (*model.Employee, error) {
	return s.EmpRepository.GetEmployeeByName(ctx, name)
}

func (s *EmployeeService) UpdateEmployeePresentAt(ctx context.Context, empID int64, presentAt time.Time) error {
	return s.EmpRepository.UpdateEmployeePresentAt(ctx, empID, presentAt)
}
