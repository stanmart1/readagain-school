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

	var user models.User
	if err := database.DB.Preload("Role.Permissions").Where("email = ?", "admin@readagain.com").First(&user).Error; err != nil {
		log.Fatal("Failed to find admin user:", err)
	}

	log.Printf("User: %s (ID: %d)", user.Email, user.ID)
	log.Printf("Role ID: %d", user.RoleID)
	
	if user.Role != nil {
		log.Printf("Role: %s (ID: %d)", user.Role.Name, user.Role.ID)
		log.Printf("Permissions count: %d", len(user.Role.Permissions))
		for _, perm := range user.Role.Permissions {
			log.Printf("  - %s", perm.Name)
		}
	} else {
		log.Println("No role assigned!")
	}
}
