package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Role struct {
	ID   uint
	Name string
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	var roles []Role
	db.Raw("SELECT id, name FROM roles ORDER BY name, id").Scan(&roles)

	fmt.Println("\n=== ALL ROLES ===")
	for _, r := range roles {
		fmt.Printf("ID: %d | Name: %s\n", r.ID, r.Name)
	}

	// Check for duplicates
	var duplicates []struct {
		Name  string
		Count int64
	}
	db.Raw("SELECT name, COUNT(*) as count FROM roles GROUP BY name HAVING COUNT(*) > 1").Scan(&duplicates)

	if len(duplicates) > 0 {
		fmt.Println("\n⚠️  DUPLICATE ROLES FOUND:")
		for _, d := range duplicates {
			fmt.Printf("  - %s: %d occurrences\n", d.Name, d.Count)
		}
	} else {
		fmt.Println("\n✓ No duplicate roles found")
	}
}
