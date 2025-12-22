package main

import (
	"fmt"
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

	db := database.DB

	// Get student role
	var studentRole models.Role
	db.Where("name = ?", "student").First(&studentRole)
	if studentRole.ID == 0 {
		log.Fatal("Student role not found")
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	students := []models.User{
		{
			Email:           "student1@school.com",
			Username:        "student1",
			PasswordHash:    string(hashedPassword),
			FirstName:       "John",
			LastName:        "Doe",
			SchoolName:      "Springfield Elementary",
			ClassLevel:      "Grade 5",
			RoleID:          studentRole.ID,
			IsActive:        true,
			IsEmailVerified: true,
		},
		{
			Email:           "student2@school.com",
			Username:        "student2",
			PasswordHash:    string(hashedPassword),
			FirstName:       "Jane",
			LastName:        "Smith",
			SchoolName:      "Springfield Elementary",
			ClassLevel:      "Grade 4",
			RoleID:          studentRole.ID,
			IsActive:        true,
			IsEmailVerified: true,
		},
		{
			Email:           "student3@school.com",
			Username:        "student3",
			PasswordHash:    string(hashedPassword),
			FirstName:       "Mike",
			LastName:        "Johnson",
			SchoolName:      "Riverside School",
			ClassLevel:      "Grade 6",
			RoleID:          studentRole.ID,
			IsActive:        true,
			IsEmailVerified: true,
		},
		{
			Email:           "student4@school.com",
			Username:        "student4",
			PasswordHash:    string(hashedPassword),
			FirstName:       "Sarah",
			LastName:        "Williams",
			SchoolName:      "Riverside School",
			ClassLevel:      "Grade 5",
			RoleID:          studentRole.ID,
			IsActive:        true,
			IsEmailVerified: true,
		},
		{
			Email:           "student5@school.com",
			Username:        "student5",
			PasswordHash:    string(hashedPassword),
			FirstName:       "Tom",
			LastName:        "Brown",
			SchoolName:      "Springfield Elementary",
			ClassLevel:      "Grade 3",
			RoleID:          studentRole.ID,
			IsActive:        true,
			IsEmailVerified: true,
		},
	}

	for _, student := range students {
		if err := db.Create(&student).Error; err != nil {
			log.Printf("Failed to create student %s: %v", student.Email, err)
			continue
		}
		fmt.Printf("Created student: %s %s (%s)\n", student.FirstName, student.LastName, student.Email)
	}

	fmt.Println("\n✅ Successfully seeded students!")
}
