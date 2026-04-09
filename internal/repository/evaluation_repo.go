package repository

import (
	"context"
	"kpi-backend/internal/domain"

	"gorm.io/gorm"
)

type evaluationRepository struct {
	db *gorm.DB
}

func NewEvaluationRepository(db *gorm.DB) domain.EvaluationRepository {
	return &evaluationRepository{db: db}
}

func (r *evaluationRepository) Create(ctx context.Context, evaluation *domain.EvaluationResult) error {
	return r.db.WithContext(ctx).Create(evaluation).Error
}

func (r *evaluationRepository) List(ctx context.Context, filter domain.EvaluationFilter) ([]domain.EvaluationResponseDTO, error) {
	var results []domain.EvaluationResponseDTO
	query := r.db.WithContext(ctx).Model(&domain.EvaluationResult{}).Where("evaluation_results.tenant_id = ?", filter.TenantID)
	if filter.SummaryOnly {
		query = query.Omit("details")
	}

	// Always JOIN employee + department to populate virtual fields
	query = query.
		Joins("LEFT JOIN employees ON employees.id = evaluation_results.employee_id").
		Joins("LEFT JOIN departments ON departments.id = employees.department_id").
		Joins("LEFT JOIN users ON users.id = evaluation_results.evaluator_id").
		Select(`evaluation_results.*, 
			COALESCE(employees.first_name || ' ' || employees.last_name, evaluation_results.employee_name) AS employee_name,
			COALESCE(users.full_name, evaluation_results.evaluator_name) AS evaluator_name,
			COALESCE(departments.name, '') AS department_name`)

	if filter.ViewerRole == string(domain.RoleSuperAdmin) || filter.ViewerRole == string(domain.RoleCEO) {
		// CEO and SuperAdmin see all evaluations
	} else if filter.ViewerRole == string(domain.RoleEmployee) {
		// Employees see evaluations where they are the subject
		query = query.Where("evaluation_results.employee_id = ?", filter.SubjectEmployeeID)
	} else {
		// HR, Head of Dept, and others see ONLY evaluations they performed
		query = query.Where("evaluation_results.evaluator_id = ?", filter.ViewerID)
	}

    // Advanced search params
    if filter.Period != "" {
        query = query.Where("evaluation_results.period = ?", filter.Period)
    }

    if filter.SearchQuery != "" {
        q := "%" + filter.SearchQuery + "%"
        query = query.Where("evaluation_results.employee_name ILIKE ? OR evaluation_results.evaluator_name ILIKE ?", q, q)
    }
    
    // Department filter
    if filter.DepartmentID != "" {
        query = query.Where("employees.department_id = ?", filter.DepartmentID)
    }

	err := query.Order("evaluation_results.created_at desc").Scan(&results).Error
	return results, err
}

func (r *evaluationRepository) Delete(ctx context.Context, id string, tenantID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.EvaluationResult{}).Error
}
