package usecase

import (
	"context"
	"kpi-backend/internal/domain"
)

type salaryFormulaUseCase struct {
	repo domain.SalaryFormulaRepository
}

func NewSalaryFormulaUseCase(repo domain.SalaryFormulaRepository) domain.SalaryFormulaUseCase {
	return &salaryFormulaUseCase{repo: repo}
}

func (u *salaryFormulaUseCase) GetFormula(ctx context.Context, tenantID string) (*domain.SalaryFormula, error) {
	return u.repo.GetByTenant(ctx, tenantID)
}

func (u *salaryFormulaUseCase) UpdateFormula(ctx context.Context, formula domain.SalaryFormula) (*domain.SalaryFormula, error) {
	if err := u.repo.Upsert(ctx, &formula); err != nil {
		return nil, err
	}
	return u.repo.GetByTenant(ctx, formula.TenantID)
}
