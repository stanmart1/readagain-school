package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/services"
	"readagain/internal/utils"
)

type ActivityHandler struct {
	service *services.ActivityService
}

func NewActivityHandler(service *services.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

func (h *ActivityHandler) GetActivities(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	activities, hasMore, err := h.service.GetUserActivities(userID, limit, offset)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get activities: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch activities"})
	}

	return c.JSON(fiber.Map{
		"activities": activities,
		"has_more":   hasMore,
	})
}
