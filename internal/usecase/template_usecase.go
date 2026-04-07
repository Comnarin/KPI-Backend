package usecase

import (
	"context"
	"kpi-backend/internal/domain"
)

type templateUseCase struct {
	templateRepo domain.TemplateRepository
}

func NewTemplateUseCase(templateRepo domain.TemplateRepository) domain.TemplateUseCase {
	return &templateUseCase{templateRepo: templateRepo}
}

func (u *templateUseCase) ListTemplates(ctx context.Context, filter domain.TemplateFilter) ([]domain.EvaluationTemplate, error) {
	return u.templateRepo.List(ctx, filter)
}

func (u *templateUseCase) GetTemplate(ctx context.Context, id, tenantID string) (*domain.EvaluationTemplate, error) {
	return u.templateRepo.GetByID(ctx, id)
}

func (u *templateUseCase) CreateTemplate(ctx context.Context, tpl domain.EvaluationTemplate) (*domain.EvaluationTemplate, error) {
	if err := u.templateRepo.Create(ctx, &tpl); err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (u *templateUseCase) UpdateTemplate(ctx context.Context, tpl domain.EvaluationTemplate) (*domain.EvaluationTemplate, error) {
	existing, err := u.templateRepo.GetByID(ctx, tpl.ID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	existing.Name = tpl.Name
	existing.Department = tpl.Department
	existing.Period = tpl.Period
	existing.Definition = tpl.Definition
	if err := u.templateRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (u *templateUseCase) DeleteTemplate(ctx context.Context, id string, tenantID string) error {
	return u.templateRepo.Delete(ctx, id, tenantID)
}
