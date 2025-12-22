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

	log.Println("✅ Database connected")

	if err := database.DB.Exec("ALTER TABLE reviews ADD COLUMN IF NOT EXISTS is_featured BOOLEAN DEFAULT false").Error; err != nil {
		log.Fatal("❌ Failed to add is_featured column:", err)
	}

	log.Println("✅ is_featured column added to reviews table!")
}
