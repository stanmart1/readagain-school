package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	roles := []struct {
		Name        string
		Description string
	}{
		{"platform_admin", "Platform Administrator"},
		{"school_admin", "School Administrator"},
		{"teacher", "Teacher"},
		{"student", "Student"},
	}

	for _, r := range roles {
		db.Exec("INSERT INTO roles (name, description) VALUES (?, ?) ON CONFLICT (name) DO NOTHING", r.Name, r.Description)
		fmt.Printf("✓ Created role: %s\n", r.Name)
	}

	fmt.Println("\n✓ All roles created")
}
