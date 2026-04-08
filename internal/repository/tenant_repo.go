package repository

import (
	"context"
	"kpi-backend/internal/domain"

	"gorm.io/gorm"
)

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) domain.TenantRepository {
	return &tenantRepository{db: db}
}

func (r *tenantRepository) Create(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Create(tenant).Error
}

func (r *tenantRepository) GetByID(ctx context.Context, id string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).First(&tenant, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) GetByCode(ctx context.Context, code string) (*domain.Tenant, error) {
	var tenant domain.Tenant
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *tenantRepository) List(ctx context.Context) ([]domain.Tenant, error) {
	var tenants []domain.Tenant
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&tenants).Error
	return tenants, err
}

func (r *tenantRepository) Update(ctx context.Context, tenant *domain.Tenant) error {
	return r.db.WithContext(ctx).Save(tenant).Error
}

func (r *tenantRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete all dependent records
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.EvaluationResult{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.EvaluationTemplate{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.SalaryAdjustment{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.SalaryFormula{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.Employee{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.User{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.Department{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.EvaluationPeriod{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.RolePermission{}).Error; err != nil { return err }
		if err := tx.Where("tenant_id = ?", id).Delete(&domain.TenantConfig{}).Error; err != nil { return err }
		
		// Finally delete the tenant
		return tx.Where("id = ?", id).Delete(&domain.Tenant{}).Error
	})
}

func (r *tenantRepository) GetConfig(ctx context.Context, tenantID string) (*domain.TenantConfig, error) {
	var config domain.TenantConfig
	if err := r.db.WithContext(ctx).First(&config, "tenant_id = ?", tenantID).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *tenantRepository) UpdateConfig(ctx context.Context, config *domain.TenantConfig) error {
	return r.db.WithContext(ctx).Save(config).Error
}
