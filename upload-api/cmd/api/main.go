package main

import (
	"log"
	"os"

	"readagain/upload-api/internal/handlers"
	"readagain/upload-api/internal/middleware"

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
	if err := os.MkdirAll(storagePath+"/thumbnails", 0755); err != nil {
		log.Fatal("Failed to create thumbnails directory:", err)
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024, // 500MB
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigin,
		AllowMethods: "GET,POST,DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	// Initialize handler
	uploadHandler := handlers.NewUploadHandler(storagePath)

	// Routes
	api := app.Group("/api")
	
	// Upload routes
	api.Post("/upload/cover", middleware.ValidateImageUpload(), uploadHandler.UploadCover)
	api.Post("/upload/book", middleware.ValidateBookUpload(), uploadHandler.UploadBook)
	
	// File serving and deletion
	api.Get("/files/:filename", uploadHandler.ServeFile)
	api.Delete("/files/:filename", uploadHandler.DeleteFile)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

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
