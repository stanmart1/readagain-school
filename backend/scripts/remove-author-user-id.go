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

	// Drop the user_id column and its index
	db.Exec("ALTER TABLE authors DROP CONSTRAINT IF EXISTS fk_authors_user")
	db.Exec("DROP INDEX IF EXISTS idx_authors_user_id")
	db.Exec("ALTER TABLE authors DROP COLUMN IF EXISTS user_id")

	fmt.Println("✓ Removed user_id from authors table")
}
