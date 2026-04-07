package usecase

import (
	"context"
	"kpi-backend/internal/domain"
	"time"
)

type adjustmentUseCase struct {
	repo domain.AdjustmentRepository
}

func NewAdjustmentUseCase(repo domain.AdjustmentRepository) domain.AdjustmentUseCase {
	return &adjustmentUseCase{repo: repo}
}

func (u *adjustmentUseCase) CreateAdjustment(ctx context.Context, adj domain.SalaryAdjustment) (*domain.SalaryAdjustment, error) {
	adj.CreatedAt = time.Now()
	if err := u.repo.Create(ctx, &adj); err != nil {
		return nil, err
	}
	return &adj, nil
}

func (u *adjustmentUseCase) ListAdjustments(ctx context.Context, tenantID string) ([]domain.SalaryAdjustment, error) {
	return u.repo.List(ctx, tenantID)
}
