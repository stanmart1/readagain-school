package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
	"readagain/internal/models"
)

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	// Get role IDs
	var platformAdminRoleID, schoolAdminRoleID uint
	db.Raw("SELECT id FROM roles WHERE name = 'platform_admin'").Scan(&platformAdminRoleID)
	db.Raw("SELECT id FROM roles WHERE name = 'school_admin'").Scan(&schoolAdminRoleID)

	// Create platform admin
	admin := models.User{
		Email:           "admin@readagain.com",
		Username:        "admin",
		PasswordHash:    hashPassword("admin123"),
		FirstName:       "Platform",
		LastName:        "Admin",
		RoleID:          platformAdminRoleID,
		IsActive:        true,
		IsEmailVerified: true,
	}
	db.Create(&admin)
	fmt.Println("✓ Created platform admin: admin@readagain.com")

	// Create school admin
	dewayne := models.User{
		Email:           "dewayne.frazier@aun.edu.ng",
		Username:        "DeWayne",
		PasswordHash:    hashPassword("dewayne123"),
		FirstName:       "DeWayne",
		LastName:        "Frazier",
		RoleID:          schoolAdminRoleID,
		IsActive:        true,
		IsEmailVerified: true,
	}
	db.Create(&dewayne)
	fmt.Println("✓ Created school admin: dewayne.frazier@aun.edu.ng")
}
