package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/models"
	"readagain/internal/services"
	"readagain/internal/utils"
)

type TestimonialHandler struct {
	service *services.TestimonialService
}

func NewTestimonialHandler(service *services.TestimonialService) *TestimonialHandler {
	return &TestimonialHandler{service: service}
}

func (h *TestimonialHandler) List(c *fiber.Ctx) error {
	testimonials, err := h.service.List(true)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list testimonials: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch testimonials"})
	}

	return c.JSON(fiber.Map{"data": testimonials})
}

func (h *TestimonialHandler) AdminList(c *fiber.Ctx) error {
	testimonials, err := h.service.List(false)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list testimonials: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch testimonials"})
	}

	return c.JSON(fiber.Map{"data": testimonials})
}

func (h *TestimonialHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid testimonial ID"})
	}

	testimonial, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Testimonial not found"})
	}

	return c.JSON(fiber.Map{"data": testimonial})
}

func (h *TestimonialHandler) Create(c *fiber.Ctx) error {
	var testimonial models.Testimonial
	if err := c.BodyParser(&testimonial); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&testimonial); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.Create(&testimonial); err != nil {
		utils.ErrorLogger.Printf("Failed to create testimonial: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create testimonial"})
	}

	utils.InfoLogger.Printf("Testimonial created: %d", testimonial.ID)
	return c.Status(201).JSON(fiber.Map{"data": testimonial})
}

func (h *TestimonialHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid testimonial ID"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Update(uint(id), updates); err != nil {
		utils.ErrorLogger.Printf("Failed to update testimonial: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update testimonial"})
	}

	testimonial, _ := h.service.GetByID(uint(id))
	utils.InfoLogger.Printf("Testimonial updated: %d", id)
	return c.JSON(fiber.Map{"data": testimonial})
}

func (h *TestimonialHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid testimonial ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete testimonial: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete testimonial"})
	}

	utils.InfoLogger.Printf("Testimonial deleted: %d", id)
	return c.JSON(fiber.Map{"message": "Testimonial deleted successfully"})
}
