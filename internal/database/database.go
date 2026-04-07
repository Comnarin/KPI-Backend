package database

import (
	"kpi-backend/internal/domain"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectAndMigrate() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://kpi_user:kpi_password@localhost:5432/kpi?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := db.AutoMigrate(
		&domain.Tenant{},
		&domain.TenantConfig{},
		&domain.User{},
		&domain.Employee{},
		&domain.EvaluationTemplate{},
		&domain.EvaluationResult{},
		// New tables
		&domain.Department{},
		&domain.SalaryFormula{},
		&domain.RolePermission{},
		&domain.EvaluationPeriod{},
		&domain.SalaryAdjustment{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	return db
}
