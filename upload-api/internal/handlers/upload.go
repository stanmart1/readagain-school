package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"readagain/upload-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UploadHandler struct {
	storagePath string
	optimizer   *utils.ImageOptimizer
}

func NewUploadHandler(storagePath string) *UploadHandler {
	return &UploadHandler{
		storagePath: storagePath,
		optimizer:   utils.NewImageOptimizer(1200, 1800, 85),
	}
}

func (h *UploadHandler) UploadCover(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(h.storagePath, "covers", filename)

	// Save file
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	// Optimize image
	if err := h.optimizer.OptimizeImage(filePath); err != nil {
		os.Remove(filePath)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to optimize image",
		})
	}

	// Update filename if converted to jpg
	if ext == ".png" || ext == ".webp" {
		filename = strings.TrimSuffix(filename, ext) + ".jpg"
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"filename": filename,
		"path":     fmt.Sprintf("covers/%s", filename),
		"url":      fmt.Sprintf("/api/files/%s", filename),
		"size":     file.Size,
	})
}

func (h *UploadHandler) UploadBook(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(h.storagePath, "books", filename)

	// Save file
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"filename": filename,
		"url":      fmt.Sprintf("/api/files/%s", filename),
		"size":     file.Size,
	})
}

func (h *UploadHandler) ServeFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	
	// Check in covers directory
	coverPath := filepath.Join(h.storagePath, "covers", filename)
	if _, err := os.Stat(coverPath); err == nil {
		return c.SendFile(coverPath)
	}

	// Check in books directory
	bookPath := filepath.Join(h.storagePath, "books", filename)
	if _, err := os.Stat(bookPath); err == nil {
		return c.SendFile(bookPath)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "File not found",
	})
}

func (h *UploadHandler) DeleteFile(c *fiber.Ctx) error {
	filename := c.Params("filename")
	
	// Try to delete from covers
	coverPath := filepath.Join(h.storagePath, "covers", filename)
	if err := os.Remove(coverPath); err == nil {
		return c.JSON(fiber.Map{
			"message": "File deleted successfully",
		})
	}

	// Try to delete from books
	bookPath := filepath.Join(h.storagePath, "books", filename)
	if err := os.Remove(bookPath); err == nil {
		return c.JSON(fiber.Map{
			"message": "File deleted successfully",
		})
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "File not found",
	})
}
