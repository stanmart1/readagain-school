package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	
	"readagain/internal/models"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.AuthLog{},
		&models.TokenBlacklist{},
		&models.Author{},
		&models.Book{},
		&models.Category{},
		&models.UserLibrary{},
		&models.ReadingSession{},
		&models.ReadingGoal{},
		&models.Bookmark{},
		&models.Note{},
		&models.Highlight{},
		&models.Blog{},
		&models.FAQ{},
		&models.Review{},
		&models.Testimonial{},
		&models.ContactMessage{},
		&models.SystemSettings{},
		&models.AuditLog{},
		&models.Notification{},
		&models.Achievement{},
		&models.UserAchievement{},
		&models.AboutPage{},
		&models.Activity{},
		&models.Wishlist{},
		&models.Group{},
		&models.GroupMember{},
		&models.ChatRoom{},
		&models.ChatMessage{},
		&models.ChatMember{},
		&models.ChatReaction{},
		&models.TypingIndicator{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	fmt.Println("✓ AutoMigrate completed")
}
