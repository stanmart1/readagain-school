package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/services"
	"readagain/internal/utils"
)

type AuditHandler struct {
	service *services.AuditService
}

func NewAuditHandler(service *services.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) GetLogs(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	entityType := c.Query("entity_type")
	action := c.Query("action")

	var userID *uint
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if id, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			uid := uint(id)
			userID = &uid
		}
	}

	logs, meta, err := h.service.GetLogs(page, limit, userID, entityType, action)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get audit logs: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch audit logs"})
	}

	return c.JSON(fiber.Map{"data": logs, "meta": meta})
}

func (h *AuditHandler) GetLog(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid audit log ID"})
	}

	log, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Audit log not found"})
	}

	return c.JSON(fiber.Map{"data": log})
}
