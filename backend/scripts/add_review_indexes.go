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
		"CREATE INDEX IF NOT EXISTS idx_reviews_is_featured ON reviews(is_featured)",
		"CREATE INDEX IF NOT EXISTS idx_reviews_status ON reviews(status)",
		"CREATE INDEX IF NOT EXISTS idx_reviews_featured_status ON reviews(is_featured, status)",
	}

	for _, idx := range indexes {
		if err := database.DB.Exec(idx).Error; err != nil {
			log.Printf("❌ Failed to create index: %v", err)
		}
	}

	log.Println("✅ Review indexes created successfully!")
}
