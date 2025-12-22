package handlers

import (
	"github.com/gofiber/fiber/v2"

	"readagain/internal/services"
	"readagain/internal/utils"
)

type AboutHandler struct {
	service *services.AboutService
}

func NewAboutHandler(service *services.AboutService) *AboutHandler {
	return &AboutHandler{service: service}
}

func (h *AboutHandler) Get(c *fiber.Ctx) error {
	about, err := h.service.Get()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get about page: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch about page"})
	}

	return c.JSON(about)
}

func (h *AboutHandler) Update(c *fiber.Ctx) error {
	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Update(updates); err != nil {
		utils.ErrorLogger.Printf("Failed to update about page: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update about page"})
	}

	about, _ := h.service.Get()
	utils.InfoLogger.Println("About page updated")
	return c.JSON(fiber.Map{"data": about})
}
