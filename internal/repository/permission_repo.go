package repository

import (
	"context"
	"errors"
	"kpi-backend/internal/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Default permissions per role
var defaultPermissions = map[string]map[string]bool{
	string(domain.RoleCEO): {
		domain.PermViewDashboard:  true,
		domain.PermViewSalary:     true,
		domain.PermApproveSalary:  true,
		domain.PermCRUDEmployee:   true,
		domain.PermCRUDDepartment: true,
		domain.PermCRUDTemplate:   true,
		domain.PermEvalKPI:        true,
		domain.PermManageUsers:    true,
		domain.PermManageSettings: true,
		domain.PermManagePeriod:   true,
	},
	string(domain.RoleHR): {
		domain.PermViewDashboard:  true,
		domain.PermViewSalary:     false,
		domain.PermApproveSalary:  false,
		domain.PermCRUDEmployee:   true,
		domain.PermCRUDDepartment: true,
		domain.PermCRUDTemplate:   true,
		domain.PermEvalKPI:        true,
		domain.PermManageUsers:    false,
		domain.PermManageSettings: false,
		domain.PermManagePeriod:   false,
	},
	string(domain.RoleHead): {
		domain.PermViewDashboard:  true,
		domain.PermViewSalary:     true,
		domain.PermApproveSalary:  false,
		domain.PermCRUDEmployee:   false,
		domain.PermCRUDDepartment: false,
		domain.PermCRUDTemplate:   true,
		domain.PermEvalKPI:        true,
		domain.PermManageUsers:    false,
		domain.PermManageSettings: false,
		domain.PermManagePeriod:   false,
	},
	string(domain.RoleEmployee): {
		domain.PermViewDashboard:  true,
		domain.PermViewSalary:     false,
		domain.PermApproveSalary:  false,
		domain.PermCRUDEmployee:   false,
		domain.PermCRUDDepartment: false,
		domain.PermCRUDTemplate:   false,
		domain.PermEvalKPI:        false,
		domain.PermManageUsers:    false,
		domain.PermManageSettings: false,
		domain.PermManagePeriod:   false,
	},
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domain.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) InitDefaults(ctx context.Context, tenantID string) error {
	var perms []domain.RolePermission
	for role, permMap := range defaultPermissions {
		for permKey, allowed := range permMap {
			perms = append(perms, domain.RolePermission{
				TenantID:      tenantID,
				Role:          role,
				PermissionKey: permKey,
				Allowed:       allowed,
			})
		}
	}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&perms).Error
}

func (r *permissionRepository) GetForTenant(ctx context.Context, tenantID string) ([]domain.RolePermission, error) {
	var perms []domain.RolePermission
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Find(&perms).Error
	return perms, err
}

func (r *permissionRepository) BulkUpsert(ctx context.Context, perms []domain.RolePermission) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tenant_id"}, {Name: "role"}, {Name: "permission_key"}},
			DoUpdates: clause.AssignmentColumns([]string{"allowed", "updated_at"}),
		}).
		Create(&perms).Error
}

func (r *permissionRepository) HasPermission(ctx context.Context, tenantID string, role string, permKey string) (bool, error) {
	var perm domain.RolePermission
	// Take() skips the implicit ORDER BY pk that First() adds — cleaner on indexed lookup
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND role = ? AND permission_key = ?", tenantID, role, permKey).
		Take(&perm).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return perm.Allowed, nil
}
