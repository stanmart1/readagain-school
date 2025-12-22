package middleware

import (
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
)

const (
	MaxImageSize = 10 * 1024 * 1024  // 10MB
	MaxBookSize  = 100 * 1024 * 1024 // 100MB
)

var (
	AllowedImageTypes = []string{"image/jpeg", "image/jpg", "image/png", "image/webp"}
	AllowedBookTypes  = []string{"application/pdf", "application/epub+zip"}
)

func ValidateImageUpload() fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No file provided",
			})
		}

		if file.Size > MaxImageSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Image size exceeds 10MB limit",
			})
		}

		if !isAllowedType(file, AllowedImageTypes) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid image type. Allowed: JPEG, PNG, WebP",
			})
		}

		return c.Next()
	}
}

func ValidateBookUpload() fiber.Handler {
	return func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No file provided",
			})
		}

		if file.Size > MaxBookSize {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Book file size exceeds 100MB limit",
			})
		}

		if !isAllowedType(file, AllowedBookTypes) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid book type. Allowed: PDF, EPUB",
			})
		}

		return c.Next()
	}
}

func isAllowedType(file *multipart.FileHeader, allowedTypes []string) bool {
	contentType := file.Header.Get("Content-Type")
	for _, allowed := range allowedTypes {
		if strings.EqualFold(contentType, allowed) {
			return true
		}
	}
	return false
}
