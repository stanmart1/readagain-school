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

	achievements := []models.Achievement{
		{Name: "First Purchase", Description: "Purchase your first book", Icon: "📚", Type: "books_purchased", Target: 1, Points: 10},
		{Name: "Book Collector", Description: "Purchase 10 books", Icon: "📖", Type: "books_purchased", Target: 10, Points: 50},
		{Name: "Library Master", Description: "Purchase 50 books", Icon: "🏛️", Type: "books_purchased", Target: 50, Points: 200},
		{Name: "First Finish", Description: "Complete your first book", Icon: "🎉", Type: "books_completed", Target: 1, Points: 20},
		{Name: "Bookworm", Description: "Complete 5 books", Icon: "🐛", Type: "books_completed", Target: 5, Points: 100},
		{Name: "Avid Reader", Description: "Complete 25 books", Icon: "📕", Type: "books_completed", Target: 25, Points: 300},
		{Name: "Reading Legend", Description: "Complete 100 books", Icon: "👑", Type: "books_completed", Target: 100, Points: 1000},
		{Name: "Getting Started", Description: "Read for 60 minutes", Icon: "⏱️", Type: "reading_minutes", Target: 60, Points: 15},
		{Name: "Dedicated Reader", Description: "Read for 10 hours", Icon: "⏰", Type: "reading_minutes", Target: 600, Points: 100},
		{Name: "Time Master", Description: "Read for 100 hours", Icon: "🕐", Type: "reading_minutes", Target: 6000, Points: 500},
		{Name: "First Session", Description: "Complete your first reading session", Icon: "🎯", Type: "reading_sessions", Target: 1, Points: 5},
		{Name: "Consistent Reader", Description: "Complete 50 reading sessions", Icon: "📅", Type: "reading_sessions", Target: 50, Points: 150},
		{Name: "Reading Habit", Description: "Complete 200 reading sessions", Icon: "🔥", Type: "reading_sessions", Target: 200, Points: 400},
	}

	for _, achievement := range achievements {
		var existing models.Achievement
		if err := database.DB.Where("name = ?", achievement.Name).First(&existing).Error; err != nil {
			if err := database.DB.Create(&achievement).Error; err != nil {
				log.Printf("Failed to create achievement %s: %v", achievement.Name, err)
			} else {
				log.Printf("✓ Created achievement: %s", achievement.Name)
			}
		} else {
			log.Printf("- Achievement already exists: %s", achievement.Name)
		}
	}

	log.Println("\n✅ Achievement seeding completed!")
}
