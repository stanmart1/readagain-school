package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
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

	var user models.User
	err := database.DB.Where("email = ?", "admin@readagain.com").First(&user).Error
	if err != nil {
		log.Fatal("❌ User not found:", err)
	}

	log.Printf("✅ User found:")
	log.Printf("   ID: %d", user.ID)
	log.Printf("   Email: %s", user.Email)
	log.Printf("   Username: %s", user.Username)
	log.Printf("   IsActive: %v", user.IsActive)
	log.Printf("   IsEmailVerified: %v", user.IsEmailVerified)
	log.Printf("   RoleID: %d", user.RoleID)

	password := "admin123"
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		log.Printf("❌ Password verification failed: %v", err)
	} else {
		log.Printf("✅ Password verification successful!")
	}
}
