package repository

import (
	"context"
	"kpi-backend/internal/domain"

	"gorm.io/gorm"
)

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) domain.DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) List(ctx context.Context, filter domain.DepartmentFilter) ([]domain.Department, error) {
	var depts []domain.Department
	query := r.db.WithContext(ctx).Where("tenant_id = ?", filter.TenantID)

	if filter.SearchQuery != "" {
		q := "%" + filter.SearchQuery + "%"
		query = query.Where("name ILIKE ?", q)
	}

	err := query.Order("created_at desc").Find(&depts).Error
	return depts, err
}

func (r *departmentRepository) GetByID(ctx context.Context, id string) (*domain.Department, error) {
	var result domain.Department
	if err := r.db.WithContext(ctx).First(&result, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *departmentRepository) Create(ctx context.Context, dept *domain.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

func (r *departmentRepository) Update(ctx context.Context, dept *domain.Department) error {
	return r.db.WithContext(ctx).Save(dept).Error
}

func (r *departmentRepository) Delete(ctx context.Context, id string, tenantID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.Department{}).Error
}
