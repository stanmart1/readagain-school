package main

import (
	"log"
	"os"

	"readagain/upload-api/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get configuration
	port := getEnv("PORT", "8001")
	storagePath := getEnv("COOLIFY_STORAGE_PATH", "/app/storage")
	corsOrigin := getEnv("CORS_ORIGIN", "*")

	// Create storage directories
	if err := os.MkdirAll(storagePath+"/covers", 0755); err != nil {
		log.Fatal("Failed to create covers directory:", err)
	}
	if err := os.MkdirAll(storagePath+"/books", 0755); err != nil {
		log.Fatal("Failed to create books directory:", err)
	}
	if err := os.MkdirAll(storagePath+"/profiles", 0755); err != nil {
		log.Fatal("Failed to create profiles directory:", err)
	}
	if err := os.MkdirAll(storagePath+"/thumbnails", 0755); err != nil {
		log.Fatal("Failed to create thumbnails directory:", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit:             500 * 1024 * 1024, // 500MB
		StreamRequestBody:     true,               // Enable streaming for large files
		DisableStartupMessage: false,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigin,
		AllowMethods: "GET,POST,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Setup routes
	routes.SetupRoutes(app, storagePath)

	// Start server
	log.Printf("🚀 Upload API starting on port %s", port)
	log.Printf("📁 Storage path: %s", storagePath)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
