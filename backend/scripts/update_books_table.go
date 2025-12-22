package main

import (
	"log"
	"readagain/internal/config"
	"readagain/internal/database"
)

func main() {
	cfg := config.Load()
	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Drop price column
	if err := database.DB.Exec("ALTER TABLE books DROP COLUMN IF EXISTS price").Error; err != nil {
		log.Printf("Failed to drop price column: %v", err)
	} else {
		log.Println("✓ Dropped price column from books table")
	}

	// Rename download_count to library_count
	if err := database.DB.Exec("ALTER TABLE books RENAME COLUMN download_count TO library_count").Error; err != nil {
		log.Printf("Failed to rename download_count column: %v", err)
	} else {
		log.Println("✓ Renamed download_count to library_count in books table")
	}

	log.Println("\n✅ Books table updated successfully!")
}
