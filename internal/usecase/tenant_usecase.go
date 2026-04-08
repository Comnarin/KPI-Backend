package usecase

import (
	"context"
	"kpi-backend/internal/domain"
	"time"
)

type tenantUseCase struct {
	tenantRepo domain.TenantRepository
	periodRepo domain.PeriodRepository
	userRepo   domain.UserRepository
	permRepo   domain.PermissionRepository
}

func NewTenantUseCase(
	tenantRepo domain.TenantRepository,
	periodRepo domain.PeriodRepository,
	userRepo domain.UserRepository,
	permRepo domain.PermissionRepository,
) domain.TenantUseCase {
	return &tenantUseCase{
		tenantRepo: tenantRepo,
		periodRepo: periodRepo,
		userRepo:   userRepo,
		permRepo:   permRepo,
	}
}

func (u *tenantUseCase) CreateTenant(ctx context.Context, tenant domain.Tenant, adminEmail, adminPassword, adminName string) (*domain.Tenant, error) {
	tenant.Config.EnableCEOEval = true // Default
	if err := u.tenantRepo.Create(ctx, &tenant); err != nil {
		return nil, err
	}

	// Create first admin user
	if adminEmail != "" {
		admin := &domain.User{
			TenantID: tenant.ID,
			Email:    adminEmail,
			Password: adminPassword,
			FullName: adminName,
			Role:     domain.RoleCEO,
		}
		if err := u.userRepo.Create(ctx, admin); err != nil {
			return nil, err
		}
	}

	// Initialize default role permissions
	if err := u.permRepo.InitDefaults(ctx, tenant.ID); err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (u *tenantUseCase) GetTenant(ctx context.Context, id string) (*domain.Tenant, error) {
	return u.tenantRepo.GetByID(ctx, id)
}

func (u *tenantUseCase) ListTenants(ctx context.Context) ([]domain.Tenant, error) {
	return u.tenantRepo.List(ctx)
}

func (u *tenantUseCase) UpdateTenant(ctx context.Context, id string, name, code, size string, maxUsers int, isActive bool) (*domain.Tenant, error) {
	tenant, err := u.tenantRepo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	tenant.Name = name
	tenant.Code = code
	tenant.Size = size
	tenant.MaxUsers = maxUsers
	tenant.IsActive = isActive

	if err := u.tenantRepo.Update(ctx, tenant); err != nil {
		return nil, err
	}
	return tenant, nil
}

func (u *tenantUseCase) DeleteTenant(ctx context.Context, id string) error {
	// Add business logic to check if tenant can be deleted, e.g. wiping associated data
	// For now we just delete the tenant record
	return u.tenantRepo.Delete(ctx, id)
}

func (u *tenantUseCase) GetTenantConfig(ctx context.Context, tenantID string) (*domain.TenantConfig, error) {
	return u.tenantRepo.GetConfig(ctx, tenantID)
}

func (u *tenantUseCase) UpdateTenantConfig(ctx context.Context, tenantID string, config domain.TenantConfig) (*domain.TenantConfig, error) {
	current, err := u.tenantRepo.GetConfig(ctx, tenantID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	current.EnableHREval = config.EnableHREval
	current.EnableDeptHeadEval = config.EnableDeptHeadEval
	current.EnableCEOEval = config.EnableCEOEval
	if err := u.tenantRepo.UpdateConfig(ctx, current); err != nil {
		return nil, err
	}
	return current, nil
}

// ─── Period Management ────────────────────────────────────────────────────────

func (u *tenantUseCase) GetActivePeriod(ctx context.Context, tenantID string) (*domain.EvaluationPeriod, error) {
	return u.periodRepo.GetActive(ctx, tenantID)
}

func (u *tenantUseCase) ListPeriods(ctx context.Context, tenantID string) ([]domain.EvaluationPeriod, error) {
	return u.periodRepo.ListByTenant(ctx, tenantID)
}

func (u *tenantUseCase) CreatePeriod(ctx context.Context, tenantID string, label string, startDate, endDate *time.Time) (*domain.EvaluationPeriod, error) {
	period := &domain.EvaluationPeriod{
		TenantID:  tenantID,
		Label:     label,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  false, // manually set active via SetActivePeriod
	}
	// If it's the first period, make it active
	existing, err := u.periodRepo.ListByTenant(ctx, tenantID)
	if err == nil && len(existing) == 0 {
		period.IsActive = true
	}
	
	if err := u.periodRepo.Create(ctx, period); err != nil {
		return nil, err
	}
	return period, nil
}

func (u *tenantUseCase) UpdatePeriod(ctx context.Context, tenantID, periodID string, label string, startDate, endDate *time.Time) (*domain.EvaluationPeriod, error) {
	period, err := u.periodRepo.GetByID(ctx, tenantID, periodID)
	if err != nil {
		return nil, err
	}
	period.Label = label
	period.StartDate = startDate
	period.EndDate = endDate
	if err := u.periodRepo.Update(ctx, period); err != nil {
		return nil, err
	}
	return period, nil
}

func (u *tenantUseCase) DeletePeriod(ctx context.Context, tenantID, periodID string) error {
	return u.periodRepo.Delete(ctx, tenantID, periodID)
}

func (u *tenantUseCase) SetActivePeriod(ctx context.Context, tenantID, periodID string) error {
	return u.periodRepo.SetActive(ctx, tenantID, periodID)
}

