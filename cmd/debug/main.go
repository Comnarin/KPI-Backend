package main

import (
	"fmt"
	"log"
	"kpi-backend/internal/database"
	"kpi-backend/internal/domain"
	"encoding/json"
)

func main() {
	db := database.ConnectAndMigrate()
	var employees []domain.Employee
	if err := db.Find(&employees).Error; err != nil {
		log.Fatal(err)
	}
	
	bytes, _ := json.MarshalIndent(employees, "", "  ")
	fmt.Println(string(bytes))
}
