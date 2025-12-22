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
		"CREATE INDEX IF NOT EXISTS idx_permissions_id ON permissions(id)",
		"CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name)",
		"CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id)",
		"CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id)",
	}

	for _, idx := range indexes {
		if err := database.DB.Exec(idx).Error; err != nil {
			log.Printf("❌ Failed to create index: %v", err)
		}
	}

	log.Println("✅ Permission indexes created successfully!")
}
