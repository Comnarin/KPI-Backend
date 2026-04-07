package repository

import (
	"context"
	"kpi-backend/internal/domain"

	"gorm.io/gorm"
)

type adjustmentRepository struct {
	db *gorm.DB
}

func NewAdjustmentRepository(db *gorm.DB) domain.AdjustmentRepository {
	return &adjustmentRepository{db: db}
}

func (r *adjustmentRepository) Create(ctx context.Context, adj *domain.SalaryAdjustment) error {
	return r.db.WithContext(ctx).Create(adj).Error
}

func (r *adjustmentRepository) List(ctx context.Context, tenantID string) ([]domain.SalaryAdjustment, error) {
	var results []domain.SalaryAdjustment
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at desc").
		Find(&results).Error
	return results, err
}
