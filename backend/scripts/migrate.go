package main

import (
	"log"

	"readagain/internal/config"
	"readagain/internal/database"
	"readagain/internal/utils"
	"readagain/migrations"
)

func main() {
	utils.InitLogger()
	cfg := config.Load()

	if err := database.Connect(cfg.Database.URL); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := migrations.RunMigrations(database.DB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	utils.InfoLogger.Println("✅ All migrations completed successfully")
}
