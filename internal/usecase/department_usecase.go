package usecase

import (
	"context"
	"kpi-backend/internal/domain"
)

type departmentUseCase struct {
	deptRepo     domain.DepartmentRepository
	employeeRepo domain.EmployeeRepository
}

func NewDepartmentUseCase(deptRepo domain.DepartmentRepository, employeeRepo domain.EmployeeRepository) domain.DepartmentUseCase {
	return &departmentUseCase{
		deptRepo:     deptRepo,
		employeeRepo: employeeRepo,
	}
}

func (u *departmentUseCase) ListDepartments(ctx context.Context, filter domain.DepartmentFilter) ([]domain.Department, error) {
	return u.deptRepo.List(ctx, filter)
}

func (u *departmentUseCase) CreateDepartment(ctx context.Context, dept domain.Department) (*domain.Department, error) {
	if err := u.deptRepo.Create(ctx, &dept); err != nil {
		return nil, err
	}
	return &dept, nil
}

func (u *departmentUseCase) UpdateDepartment(ctx context.Context, dept domain.Department) (*domain.Department, error) {
	existing, err := u.deptRepo.GetByID(ctx, dept.ID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	existing.Name = dept.Name
	existing.Code = dept.Code
	existing.Description = dept.Description
	if err := u.deptRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *departmentUseCase) DeleteDepartment(ctx context.Context, id string, tenantID string) error {
	// Check if any employees belong to this department
	emps, err := u.employeeRepo.List(ctx, domain.EmployeeFilter{
		TenantID:     tenantID,
		DepartmentID: id,
	})
	if err == nil && len(emps) > 0 {
		return domain.ErrHasEmployees
	}

	return u.deptRepo.Delete(ctx, id, tenantID)
}
