package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
	"readagain/internal/models"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	roles := []models.Role{
		{Name: "platform_admin", Description: "Platform Administrator"},
		{Name: "school_admin", Description: "School Administrator"},
		{Name: "teacher", Description: "Teacher"},
		{Name: "student", Description: "Student"},
	}

	for _, r := range roles {
		db.FirstOrCreate(&r, models.Role{Name: r.Name})
		fmt.Printf("✓ Created role: %s\n", r.Name)
	}

	fmt.Println("\n✓ All 4 roles created")
}
