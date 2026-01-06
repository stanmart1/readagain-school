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

	// Start transaction
	tx := db.Begin()

	// Keep only the first occurrence of each role name, delete duplicates
	result := tx.Exec(`
		DELETE FROM roles 
		WHERE id NOT IN (
			SELECT MIN(id) 
			FROM roles 
			GROUP BY name
		)
	`)

	if result.Error != nil {
		tx.Rollback()
		log.Fatal("Failed to delete duplicates:", result.Error)
	}

	fmt.Printf("✓ Deleted %d duplicate roles\n", result.RowsAffected)

	// Verify final state
	var roles []struct {
		ID   uint
		Name string
	}
	tx.Raw("SELECT id, name FROM roles ORDER BY id").Scan(&roles)

	fmt.Println("\n=== REMAINING ROLES ===")
	for _, r := range roles {
		fmt.Printf("ID: %d | Name: %s\n", r.ID, r.Name)
	}

	tx.Commit()
	fmt.Println("\n✓ Database cleaned successfully")
}
