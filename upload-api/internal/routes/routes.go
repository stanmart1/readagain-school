package routes

import (
	"readagain/upload-api/internal/handlers"
	"readagain/upload-api/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, storagePath string) {
	// Initialize handler
	uploadHandler := handlers.NewUploadHandler(storagePath)

	// API routes
	api := app.Group("/api")
	
	// Upload routes
	api.Post("/upload/cover", middleware.ValidateImageUpload(), uploadHandler.UploadCover)
	api.Post("/upload/book", middleware.ValidateBookUpload(), uploadHandler.UploadBook)
	api.Post("/upload/profile", middleware.ValidateImageUpload(), uploadHandler.UploadProfile)
	
	// File serving and deletion
	api.Get("/files/:filename", uploadHandler.ServeFile)
	api.Delete("/files/:filename", uploadHandler.DeleteFile)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
