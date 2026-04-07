package repository

import (
	"context"
	"kpi-backend/internal/domain"

	"gorm.io/gorm"
)

type periodRepository struct {
	db *gorm.DB
}

func NewPeriodRepository(db *gorm.DB) domain.PeriodRepository {
	return &periodRepository{db: db}
}

func (r *periodRepository) ListByTenant(ctx context.Context, tenantID string) ([]domain.EvaluationPeriod, error) {
	var periods []domain.EvaluationPeriod
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("created_at desc").
		Find(&periods).Error
	return periods, err
}

func (r *periodRepository) GetActive(ctx context.Context, tenantID string) (*domain.EvaluationPeriod, error) {
	var period domain.EvaluationPeriod
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_active = ?", tenantID, true).
		First(&period).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &period, nil
}

func (r *periodRepository) GetByID(ctx context.Context, tenantID, id string) (*domain.EvaluationPeriod, error) {
	var period domain.EvaluationPeriod
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&period).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &period, nil
}

func (r *periodRepository) Create(ctx context.Context, period *domain.EvaluationPeriod) error {
	return r.db.WithContext(ctx).Create(period).Error
}

func (r *periodRepository) Update(ctx context.Context, period *domain.EvaluationPeriod) error {
	return r.db.WithContext(ctx).Save(period).Error
}

func (r *periodRepository) Delete(ctx context.Context, tenantID, id string) error {
	res := r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.EvaluationPeriod{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *periodRepository) SetActive(ctx context.Context, tenantID, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Reset all to not active
		if err := tx.Model(&domain.EvaluationPeriod{}).
			Where("tenant_id = ?", tenantID).
			Update("is_active", false).Error; err != nil {
			return err
		}

		// Set the chosen one to active
		res := tx.Model(&domain.EvaluationPeriod{}).
			Where("id = ? AND tenant_id = ?", id, tenantID).
			Update("is_active", true)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return domain.ErrNotFound
		}
		return nil
	})
}
