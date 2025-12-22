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

	// Add indexes to users table
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);",
		"CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);",
		"CREATE INDEX IF NOT EXISTS idx_users_is_email_verified ON users(is_email_verified);",
		"CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);",
		"CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);",
	}

	for _, indexSQL := range indexes {
		if err := database.DB.Exec(indexSQL).Error; err != nil {
			log.Printf("❌ Failed to create index: %s - %v", indexSQL, err)
		} else {
			log.Printf("✅ Index created: %s", indexSQL)
		}
	}

	log.Println("✅ User indexes migration completed")
}
