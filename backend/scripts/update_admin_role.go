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

	var platformAdminRole models.Role
	if err := database.DB.Where("name = ?", "platform_admin").First(&platformAdminRole).Error; err != nil {
		log.Fatal("❌ platform_admin role not found. Run seed_roles.go first")
	}

	var existingUser models.User
	err := database.DB.Where("email = ?", "admin@readagain.com").First(&existingUser).Error
	if err != nil {
		log.Fatal("❌ Admin user not found. Run seed_admin_user.go first")
	}

	existingUser.RoleID = platformAdminRole.ID
	if err := database.DB.Save(&existingUser).Error; err != nil {
		log.Fatal("❌ Failed to update admin user role:", err)
	}

	log.Printf("✅ Updated admin user to platform_admin role (Role ID: %d)", platformAdminRole.ID)
}
