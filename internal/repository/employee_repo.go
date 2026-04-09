package repository

import (
	"context"
	"kpi-backend/internal/domain"

	"gorm.io/gorm"
)

type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) domain.EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) List(ctx context.Context, filter domain.EmployeeFilter) ([]domain.Employee, error) {
	var results []domain.Employee
	query := r.db.WithContext(ctx).
		Preload("Department").
		Where("employees.tenant_id = ?", filter.TenantID)
	
	if filter.Status != "" {
		query = query.Where("employees.status = ?", filter.Status)
	}

	if filter.DepartmentID != "" {
		query = query.Where("employees.department_id = ?", filter.DepartmentID)
	}

	if filter.SearchQuery != "" {
		q := "%" + filter.SearchQuery + "%"
		query = query.Joins("LEFT JOIN departments ON departments.id = employees.department_id").
			Where("employees.first_name ILIKE ? OR employees.last_name ILIKE ? OR employees.code ILIKE ? OR departments.name ILIKE ?", q, q, q, q)
	}

	if filter.ExcludeEvaluatedPeriod != "" && filter.EvaluatorID != "" {
		subQuery := r.db.Model(&domain.EvaluationResult{}).
			Select("employee_id").
			Where("period = ? AND evaluator_id = ?", filter.ExcludeEvaluatedPeriod, filter.EvaluatorID)
		query = query.Where("employees.id NOT IN (?)", subQuery)
	}

	err := query.Order("employees.code").Find(&results).Error
	return results, err
}

func (r *employeeRepository) GetByID(ctx context.Context, id string) (*domain.Employee, error) {
	var result domain.Employee
	err := r.db.WithContext(ctx).Preload("Department").First(&result, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *employeeRepository) GetByEmail(ctx context.Context, tenantID, email string) (*domain.Employee, error) {
	var emp domain.Employee
	err := r.db.WithContext(ctx).Preload("Department").Where("tenant_id = ? AND email = ?", tenantID, email).First(&emp).Error
	if err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) Create(ctx context.Context, emp *domain.Employee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}

func (r *employeeRepository) Update(ctx context.Context, emp *domain.Employee) error {
	return r.db.WithContext(ctx).Save(emp).Error
}

func (r *employeeRepository) UpdateFields(ctx context.Context, id string, tenantID string, fields map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&domain.Employee{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Updates(fields).Error
}

func (r *employeeRepository) Delete(ctx context.Context, id string, tenantID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Employee{}).Error
}
