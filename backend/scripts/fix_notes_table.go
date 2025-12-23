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

	log.Println("Starting notes table migration...")

	// AutoMigrate will add missing columns without dropping existing data
	if err := db.AutoMigrate(&models.Note{}); err != nil {
		log.Fatal("Failed to migrate notes table:", err)
	}
	log.Println("✓ Notes table migrated successfully")

	log.Println("Migration completed!")
}
