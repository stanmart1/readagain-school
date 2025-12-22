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

	// Create contact_messages table
	if err := database.DB.AutoMigrate(&models.ContactMessage{}); err != nil {
		log.Fatal("❌ Failed to create contact_messages table:", err)
	}

	log.Println("✅ contact_messages table created successfully")
}
