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

	log.Println("✅ Database connected")

	if err := database.DB.AutoMigrate(&models.AboutPage{}); err != nil {
		log.Fatal("❌ Failed to migrate about_pages table:", err)
	}

	log.Println("✅ about_pages table migration completed!")
}
