package repository

import (
	"context"
	"kpi-backend/internal/domain"

	"gorm.io/gorm"
)

type templateRepository struct {
	db *gorm.DB
}

func NewTemplateRepository(db *gorm.DB) domain.TemplateRepository {
	return &templateRepository{db: db}
}

func (r *templateRepository) List(ctx context.Context, filter domain.TemplateFilter) ([]domain.EvaluationTemplate, error) {
	var templates []domain.EvaluationTemplate
	query := r.db.WithContext(ctx).Where("tenant_id = ?", filter.TenantID)

	if filter.SearchQuery != "" {
		q := "%" + filter.SearchQuery + "%"
		query = query.Where("name ILIKE ?", q)
	}

	err := query.Order("created_at desc").Find(&templates).Error
	return templates, err
}

func (r *templateRepository) GetByID(ctx context.Context, id string) (*domain.EvaluationTemplate, error) {
	var result domain.EvaluationTemplate
	err := r.db.WithContext(ctx).First(&result, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *templateRepository) Create(ctx context.Context, tpl *domain.EvaluationTemplate) error {
	return r.db.WithContext(ctx).Create(tpl).Error
}

func (r *templateRepository) Update(ctx context.Context, tpl *domain.EvaluationTemplate) error {
	return r.db.WithContext(ctx).Save(tpl).Error
}

func (r *templateRepository) Delete(ctx context.Context, id string, tenantID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&domain.EvaluationTemplate{}).Error
}
