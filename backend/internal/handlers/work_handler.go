package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"readagain/internal/models"
	"readagain/internal/services"
	"readagain/internal/utils"
)

type WorkHandler struct {
	service *services.WorkService
}

func NewWorkHandler(service *services.WorkService) *WorkHandler {
	return &WorkHandler{service: service}
}

func (h *WorkHandler) GetAll(c *fiber.Ctx) error {
	works, err := h.service.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch works"})
	}
	return c.JSON(fiber.Map{"success": true, "works": works})
}

func (h *WorkHandler) GetAllAdmin(c *fiber.Ctx) error {
	works, err := h.service.GetAllAdmin()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch works"})
	}
	return c.JSON(fiber.Map{"data": works})
}

func (h *WorkHandler) Create(c *fiber.Ctx) error {
	var work models.Work
	if err := c.BodyParser(&work); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&work); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.Create(&work); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create work"})
	}

	return c.Status(201).JSON(fiber.Map{"data": work})
}

func (h *WorkHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid work ID"})
	}

	var work models.Work
	if err := c.BodyParser(&work); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Update(uint(id), &work); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update work"})
	}

	return c.JSON(fiber.Map{"message": "Work updated successfully"})
}

func (h *WorkHandler) Toggle(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid work ID"})
	}

	if err := h.service.Toggle(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to toggle work"})
	}

	return c.JSON(fiber.Map{"message": "Work toggled successfully"})
}

func (h *WorkHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid work ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete work"})
	}

	return c.JSON(fiber.Map{"message": "Work deleted successfully"})
}
