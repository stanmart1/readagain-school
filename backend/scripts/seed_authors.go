package main

import (
	"fmt"
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

	db := database.DB

	// Get platform admin user
	var adminUser models.User
	db.Joins("JOIN roles ON users.role_id = roles.id").
		Where("roles.name = ?", "platform_admin").
		First(&adminUser)
	
	if adminUser.ID == 0 {
		log.Fatal("No platform admin found")
	}

	// Create sample authors
	authors := []models.Author{
		{
			UserID:       adminUser.ID,
			BusinessName: "Classic Literature Press",
			Bio:          "Publisher of timeless classic literature for young readers",
			Status:       "active",
		},
		{
			UserID:       adminUser.ID,
			BusinessName: "Modern Children's Books",
			Bio:          "Contemporary stories for the modern child",
			Status:       "active",
		},
		{
			UserID:       adminUser.ID,
			BusinessName: "Educational Publishers",
			Bio:          "Quality educational content for schools",
			Status:       "active",
		},
	}

	for _, author := range authors {
		if err := db.Create(&author).Error; err != nil {
			log.Printf("Failed to create author %s: %v", author.BusinessName, err)
			continue
		}
		fmt.Printf("Created author: %s\n", author.BusinessName)
	}

	fmt.Println("\n✅ Successfully seeded authors!")
}
