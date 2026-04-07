package usecase

import (
	"context"
	"kpi-backend/internal/domain"
)

type permissionUseCase struct {
	permRepo domain.PermissionRepository
}

func NewPermissionUseCase(permRepo domain.PermissionRepository) domain.PermissionUseCase {
	return &permissionUseCase{permRepo: permRepo}
}

func (u *permissionUseCase) GetPermissions(ctx context.Context, tenantID string) ([]domain.RolePermission, error) {
	return u.permRepo.GetForTenant(ctx, tenantID)
}

func (u *permissionUseCase) UpdatePermissions(ctx context.Context, tenantID string, perms []domain.RolePermission) error {
	for i := range perms {
		perms[i].TenantID = tenantID
	}
	return u.permRepo.BulkUpsert(ctx, perms)
}

func (u *permissionUseCase) HasPermission(ctx context.Context, tenantID string, role string, permKey string) (bool, error) {
	// SuperAdmin always has all permissions
	if role == string(domain.RoleSuperAdmin) {
		return true, nil
	}
	return u.permRepo.HasPermission(ctx, tenantID, role, permKey)
}
