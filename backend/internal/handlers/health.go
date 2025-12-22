package handlers

import "github.com/gofiber/fiber/v2"

func GetRoot(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "ReadAgain API is running",
		"version": "1.0.0",
	})
}

func GetHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"message": "ReadAgain API is running",
	})
}
