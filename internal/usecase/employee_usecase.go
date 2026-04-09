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
	existing.DepartmentID = emp.DepartmentID
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

func (u *employeeUseCase) PartialUpdateEmployee(ctx context.Context, id string, tenantID string, req domain.UpdateEmployeeRequest) (*domain.Employee, error) {
	fields := make(map[string]interface{})
	if req.FirstName != nil {
		fields["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		fields["last_name"] = *req.LastName
	}
	if req.Email != nil {
		fields["email"] = *req.Email
	}
	if req.DepartmentID != nil {
		fields["department_id"] = *req.DepartmentID
	}
	if req.Position != nil {
		fields["position"] = *req.Position
	}
	if req.BaseSalary != nil {
		fields["base_salary"] = *req.BaseSalary
	}
	if req.PersonalCapacity != nil {
		fields["personal_capacity"] = *req.PersonalCapacity
	}
	if req.VariablePayBase != nil {
		fields["variable_pay_base"] = *req.VariablePayBase
	}
	if req.Code != nil {
		fields["code"] = *req.Code
	}
	if req.Status != nil {
		fields["status"] = *req.Status
	}
	if req.YearsOfService != nil {
		fields["years_of_service"] = *req.YearsOfService
	}
	if req.StartDate != nil {
		fields["start_date"] = *req.StartDate
	}

	if len(fields) == 0 {
		return u.employeeRepo.GetByID(ctx, id)
	}

	if err := u.employeeRepo.UpdateFields(ctx, id, tenantID, fields); err != nil {
		return nil, err
	}
	return u.employeeRepo.GetByID(ctx, id)
}

func (u *employeeUseCase) DeleteEmployee(ctx context.Context, id string, tenantID string) error {
	return u.employeeRepo.Delete(ctx, id, tenantID)
}
