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

	var existingUser models.User
	err := database.DB.Where("email = ?", "admin@readagain.com").First(&existingUser).Error
	if err == nil {
		log.Println("⚠️  Admin user already exists, updating password...")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		database.DB.Model(&existingUser).Updates(map[string]interface{}{
			"password_hash":    string(hashedPassword),
			"is_active":        true,
			"is_email_verified": true,
		})
		log.Printf("✅ Admin user password updated!")
		log.Printf("   Email: admin@readagain.com")
		log.Printf("   Username: admin")
		log.Printf("   Password: admin123")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
	}

	adminUser := models.User{
		FirstName:       "Admin",
		LastName:        "User",
		Email:           "admin@readagain.com",
		Username:        "admin",
		PasswordHash:    string(hashedPassword),
		RoleID:          1, // SuperAdmin role
		IsActive:        true,
		IsEmailVerified: true,
	}

	if err := database.DB.Create(&adminUser).Error; err != nil {
		log.Fatal("❌ Failed to create admin user:", err)
	}

	log.Printf("✅ Admin user created successfully!")
	log.Printf("   Email: admin@readagain.com")
	log.Printf("   Username: admin")
	log.Printf("   Password: admin123")
	log.Printf("   Role: SuperAdmin (ID: 1)")
}
