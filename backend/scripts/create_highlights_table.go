package main

import (
	"log"

	"readagain/internal/config"
	"readagain/internal/database"
	"readagain/internal/models"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db := database.DB

	log.Println("Starting highlights table migration...")

	// Create highlights table
	if err := db.AutoMigrate(&models.Highlight{}); err != nil {
		log.Fatal("Failed to create highlights table:", err)
	}
	log.Println("✓ Highlights table created")

	// Remove highlight column from notes table
	if db.Migrator().HasColumn(&models.Note{}, "highlight") {
		if err := db.Migrator().DropColumn(&models.Note{}, "highlight"); err != nil {
			log.Fatal("Failed to drop highlight column from notes:", err)
		}
		log.Println("✓ Removed highlight column from notes table")
	} else {
		log.Println("✓ Highlight column already removed from notes table")
	}

	log.Println("Migration completed successfully!")
}
