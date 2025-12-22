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

	categories := []models.Category{
		{Name: "Fiction", Description: "Fictional stories and novels", Status: "active"},
		{Name: "Adventure", Description: "Adventure and exploration stories", Status: "active"},
		{Name: "Fantasy", Description: "Fantasy and magical stories", Status: "active"},
		{Name: "Mystery", Description: "Mystery and detective stories", Status: "active"},
		{Name: "Science Fiction", Description: "Science fiction stories", Status: "active"},
		{Name: "Classic", Description: "Classic literature", Status: "active"},
	}

	for _, cat := range categories {
		if err := db.Create(&cat).Error; err != nil {
			log.Printf("Failed to create category %s: %v", cat.Name, err)
			continue
		}
		fmt.Printf("Created category: %s\n", cat.Name)
	}

	fmt.Println("\n✅ Successfully seeded categories!")
}
