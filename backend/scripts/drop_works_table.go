package main

import (
	"log"
	"readagain/internal/config"
	"readagain/internal/database"
)

func main() {
	cfg := config.Load()
	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.DB.Exec("DROP TABLE IF EXISTS works CASCADE").Error; err != nil {
		log.Fatal("Failed to drop works table:", err)
	}

	log.Println("✓ Works table dropped successfully")
}
