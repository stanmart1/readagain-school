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

	tx := db.Begin()

	// Delete ALL roles
	tx.Exec("DELETE FROM roles")
	
	// Reset sequence
	tx.Exec("ALTER SEQUENCE roles_id_seq RESTART WITH 1")

	// Recreate the 4 roles
	tx.Exec("INSERT INTO roles (name, description) VALUES ('platform_admin', 'Platform Administrator')")
	tx.Exec("INSERT INTO roles (name, description) VALUES ('school_admin', 'School Administrator')")
	tx.Exec("INSERT INTO roles (name, description) VALUES ('teacher', 'Teacher')")
	tx.Exec("INSERT INTO roles (name, description) VALUES ('student', 'Student')")

	// Verify
	var roles []struct {
		ID   uint
		Name string
	}
	tx.Raw("SELECT id, name FROM roles ORDER BY id").Scan(&roles)

	fmt.Println("\n=== RECREATED ROLES ===")
	for _, r := range roles {
		fmt.Printf("ID: %d | Name: %s\n", r.ID, r.Name)
	}

	tx.Commit()
	fmt.Println("\n✓ Roles reset successfully")
}
