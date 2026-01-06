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

	// Drop the old unique index
	db.Exec("DROP INDEX IF EXISTS idx_books_isbn")
	
	// Create a partial unique index that only applies to non-null values
	db.Exec("CREATE UNIQUE INDEX idx_books_isbn ON books(isbn) WHERE isbn IS NOT NULL AND isbn != ''")

	fmt.Println("✓ Fixed ISBN unique constraint to allow NULL values")
}
