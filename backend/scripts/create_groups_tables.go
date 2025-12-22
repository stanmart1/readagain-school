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

	if err := database.DB.AutoMigrate(&models.Group{}, &models.GroupMember{}); err != nil {
		log.Fatal("Failed to migrate groups tables:", err)
	}

	log.Println("✅ Groups tables created successfully")
}
