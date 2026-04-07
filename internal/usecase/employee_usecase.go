package usecase

import (
	"context"
	"kpi-backend/internal/domain"
)

type employeeUseCase struct {
	employeeRepo domain.EmployeeRepository
}

func NewEmployeeUseCase(employeeRepo domain.EmployeeRepository) domain.EmployeeUseCase {
	return &employeeUseCase{employeeRepo: employeeRepo}
}

func (u *employeeUseCase) ListEmployees(ctx context.Context, filter domain.EmployeeFilter) ([]domain.Employee, error) {
	return u.employeeRepo.List(ctx, filter)
}

func (u *employeeUseCase) CreateEmployee(ctx context.Context, emp domain.Employee) (*domain.Employee, error) {
	if err := u.employeeRepo.Create(ctx, &emp); err != nil {
		return nil, err
	}
	return &emp, nil
}

func (u *employeeUseCase) UpdateEmployee(ctx context.Context, emp domain.Employee) (*domain.Employee, error) {
	existing, err := u.employeeRepo.GetByID(ctx, emp.ID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	existing.FirstName = emp.FirstName
	existing.LastName = emp.LastName
	existing.Department = emp.Department
	existing.Position = emp.Position
	existing.BaseSalary = emp.BaseSalary
	existing.PersonalCapacity = emp.PersonalCapacity
	existing.VariablePayBase = emp.VariablePayBase
	existing.Code = emp.Code
	existing.Status = emp.Status
	existing.YearsOfService = emp.YearsOfService
	existing.StartDate = emp.StartDate
	if err := u.employeeRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *employeeUseCase) DeleteEmployee(ctx context.Context, id string, tenantID string) error {
	return u.employeeRepo.Delete(ctx, id, tenantID)
}
