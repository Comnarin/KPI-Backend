package main

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type EvaluationTemplate struct {
	ID     string `gorm:"primaryKey"`
	Period string
}

func main() {
	dsn := "postgres://kpi_user:kpi_password@localhost:5432/kpi?sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Update all templates with legacy period labels to 'ALL'
	result := db.Model(&EvaluationTemplate{}).
		Where("period IN ?", []string{"รายเดือน", "รายไตรมาส", "รายปี"}).
		Update("period", "ALL")

	if result.Error != nil {
		log.Fatalf("failed to migrate template periods: %v", result.Error)
	}

	log.Printf("Successfully migrated %d templates to 'ALL' period.", result.RowsAffected)
}
