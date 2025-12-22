package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/models"
	"readagain/internal/services"
	"readagain/internal/utils"
)

type FAQHandler struct {
	service *services.FAQService
}

func NewFAQHandler(service *services.FAQService) *FAQHandler {
	return &FAQHandler{service: service}
}

func (h *FAQHandler) List(c *fiber.Ctx) error {
	category := c.Query("category")
	faqs, err := h.service.List(category, true)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list FAQs: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch FAQs"})
	}

	return c.JSON(fiber.Map{"data": faqs})
}

func (h *FAQHandler) AdminList(c *fiber.Ctx) error {
	category := c.Query("category")
	faqs, err := h.service.List(category, false)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list FAQs: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch FAQs"})
	}

	return c.JSON(fiber.Map{"data": faqs})
}

func (h *FAQHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid FAQ ID"})
	}

	faq, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "FAQ not found"})
	}

	return c.JSON(fiber.Map{"data": faq})
}

func (h *FAQHandler) Create(c *fiber.Ctx) error {
	var faq models.FAQ
	if err := c.BodyParser(&faq); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&faq); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.Create(&faq); err != nil {
		utils.ErrorLogger.Printf("Failed to create FAQ: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create FAQ"})
	}

	utils.InfoLogger.Printf("FAQ created: %d", faq.ID)
	return c.Status(201).JSON(fiber.Map{"data": faq})
}

func (h *FAQHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid FAQ ID"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Update(uint(id), updates); err != nil {
		utils.ErrorLogger.Printf("Failed to update FAQ: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update FAQ"})
	}

	faq, _ := h.service.GetByID(uint(id))
	utils.InfoLogger.Printf("FAQ updated: %d", id)
	return c.JSON(fiber.Map{"data": faq})
}

func (h *FAQHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid FAQ ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete FAQ: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete FAQ"})
	}

	utils.InfoLogger.Printf("FAQ deleted: %d", id)
	return c.JSON(fiber.Map{"message": "FAQ deleted successfully"})
}

func (h *FAQHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.service.GetCategories()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get FAQ categories: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}

	return c.JSON(fiber.Map{"data": categories})
}
