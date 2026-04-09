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
		&domain.Department{},
		&domain.SalaryFormula{},
		&domain.RolePermission{},
		&domain.EvaluationPeriod{},
		&domain.SalaryAdjustment{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	if err := seedData(db); err != nil {
		log.Printf("Warning: failed to seed data: %v", err)
	}

	return db
}


func seedData(db *gorm.DB) error {
	// 1. Ensure System Tenant exists
	var systemTenant domain.Tenant
	err := db.Where("code = ?", "SYSTEM").First(&systemTenant).Error
	if err != nil {
		systemTenant = domain.Tenant{
			Name:     "System Administration",
			Code:     "SYSTEM",
			Size:     "Enterprise",
			MaxUsers: 9999,
			IsActive: true,
		}
		if err := db.Create(&systemTenant).Error; err != nil {
			return err
		}
		// Create default config for system tenant
		config := domain.TenantConfig{
			TenantID:           systemTenant.ID,
			EnableHREval:       true,
			EnableDeptHeadEval: true,
			EnableCEOEval:      true,
		}
		db.Create(&config)
	}

	// 2. Ensure SuperAdmin exists
	var admin domain.User
	err = db.Where("email = ?", "admin@kpi.com").First(&admin).Error
	if err != nil {
		admin = domain.User{
			TenantID: systemTenant.ID,
			Email:    "admin@kpi.com",
			Password: "password",
			FullName: "System Administrator",
			Role:     domain.RoleSuperAdmin,
		}
		if err := db.Create(&admin).Error; err != nil {
			return err
		}
	}

	return nil
}
