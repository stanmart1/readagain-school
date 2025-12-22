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

	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_system_settings_category ON system_settings(category)",
		"CREATE INDEX IF NOT EXISTS idx_system_settings_key ON system_settings(key)",
	}

	for _, idx := range indexes {
		if err := database.DB.Exec(idx).Error; err != nil {
			log.Printf("❌ Failed to create index: %v", err)
		}
	}

	log.Println("✅ System settings indexes created successfully!")
}
