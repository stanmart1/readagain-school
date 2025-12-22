package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/models"
	"readagain/internal/services"
	"readagain/internal/utils"
)

type ContactHandler struct {
	service *services.ContactService
}

func NewContactHandler(service *services.ContactService) *ContactHandler {
	return &ContactHandler{service: service}
}

func (h *ContactHandler) Submit(c *fiber.Ctx) error {
	var message models.ContactMessage
	if err := c.BodyParser(&message); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&message); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.Create(&message); err != nil {
		utils.ErrorLogger.Printf("Failed to create contact message: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to submit message"})
	}

	utils.InfoLogger.Printf("Contact message received from: %s", message.Email)
	return c.Status(201).JSON(fiber.Map{"message": "Message sent successfully"})
}

func (h *ContactHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	status := c.Query("status")

	messages, meta, err := h.service.List(page, limit, status)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list contact messages: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}

	return c.JSON(fiber.Map{"data": messages, "meta": meta})
}

func (h *ContactHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid message ID"})
	}

	message, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Message not found"})
	}

	return c.JSON(fiber.Map{"data": message})
}

func (h *ContactHandler) Reply(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid message ID"})
	}

	var req struct {
		Reply string `json:"reply" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.Reply(uint(id), req.Reply); err != nil {
		utils.ErrorLogger.Printf("Failed to reply to message: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send reply"})
	}

	message, _ := h.service.GetByID(uint(id))
	utils.InfoLogger.Printf("Replied to contact message: %d", id)
	return c.JSON(fiber.Map{"data": message})
}

func (h *ContactHandler) UpdateStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid message ID"})
	}

	var req struct {
		Status string `json:"status" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.UpdateStatus(uint(id), req.Status); err != nil {
		utils.ErrorLogger.Printf("Failed to update message status: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update status"})
	}

	message, _ := h.service.GetByID(uint(id))
	utils.InfoLogger.Printf("Updated contact message status: %d", id)
	return c.JSON(fiber.Map{"data": message})
}

func (h *ContactHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid message ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete message: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete message"})
	}

	utils.InfoLogger.Printf("Contact message deleted: %d", id)
	return c.JSON(fiber.Map{"message": "Message deleted successfully"})
}
