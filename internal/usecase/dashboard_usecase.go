package usecase

import (
	"context"

	"kpi-backend/internal/domain"
)

type dashboardUseCase struct {
	repo domain.DashboardRepository
}

func NewDashboardUseCase(repo domain.DashboardRepository) domain.DashboardUseCase {
	return &dashboardUseCase{repo: repo}
}

func (u *dashboardUseCase) GetDashboardStats(ctx context.Context, tenantID string, period string) (*domain.DashboardStats, error) {
	return u.repo.GetDashboardStats(ctx, tenantID, period)
}
