package repository

import (
	"context"
	"encoding/json"
	"errors"
	"kpi-backend/internal/domain"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type salaryFormulaRepository struct {
	db *gorm.DB
}

func NewSalaryFormulaRepository(db *gorm.DB) domain.SalaryFormulaRepository {
	return &salaryFormulaRepository{db: db}
}

func (r *salaryFormulaRepository) GetByTenant(ctx context.Context, tenantID string) (*domain.SalaryFormula, error) {
	var result domain.SalaryFormula
	// Take() avoids the implicit ORDER BY primary_key that First() adds
	err := r.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Take(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Return sensible defaults when no formula configured yet
		defaultMerit, _ := json.Marshal(map[string]int{
			"ดีเยี่ยม":     8,
			"เกินเป้า":     6,
			"ได้เป้า":      4,
			"ต้องปรับปรุง": 2,
			"ไม่ผ่านเกณฑ์": 0,
		})
		return &domain.SalaryFormula{
			TenantID:          tenantID,
			TenureRatePerYear: 0.5,
			MaxTenureBonus:    10,
			MeritMap:          datatypes.JSON(defaultMerit),
		}, nil
	}
	return &result, err
}

func (r *salaryFormulaRepository) Upsert(ctx context.Context, formula *domain.SalaryFormula) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tenant_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"tenure_rate_per_year", "max_tenure_bonus", "merit_map", "updated_at"}),
		}).
		Create(formula).Error
}
