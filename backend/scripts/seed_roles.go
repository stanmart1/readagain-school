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

	roles := []models.Role{
		{Name: "student", Description: "Student who reads books"},
		{Name: "teacher", Description: "Teacher who assigns books and tracks student progress"},
		{Name: "school_admin", Description: "School administrator who manages school users and library"},
		{Name: "platform_admin", Description: "Platform administrator with full system access"},
	}

	for _, role := range roles {
		var existing models.Role
		err := database.DB.Where("name = ?", role.Name).First(&existing).Error
		if err != nil {
			if err := database.DB.Create(&role).Error; err != nil {
				log.Printf("❌ Failed to create role %s: %v", role.Name, err)
			} else {
				log.Printf("✓ Created role: %s (ID: %d)", role.Name, role.ID)
			}
		} else {
			log.Printf("- Role already exists: %s (ID: %d)", existing.Name, existing.ID)
		}
	}

	log.Println("\n✅ Role seeding completed!")
}
