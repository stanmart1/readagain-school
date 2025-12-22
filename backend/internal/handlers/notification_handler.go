package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/services"
	"readagain/internal/utils"
)

type NotificationHandler struct {
	service *services.NotificationService
}

func NewNotificationHandler(service *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	unreadOnly := c.Query("unread_only") == "true"

	notifications, meta, err := h.service.GetUserNotifications(userID, page, limit, unreadOnly)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get notifications: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}

	return c.JSON(fiber.Map{"data": notifications, "meta": meta})
}

func (h *NotificationHandler) GetNotification(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	notification, err := h.service.GetByID(uint(id), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Notification not found"})
	}

	return c.JSON(fiber.Map{"data": notification})
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	if err := h.service.MarkAsRead(uint(id), userID); err != nil {
		utils.ErrorLogger.Printf("Failed to mark notification as read: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	if err := h.service.MarkAllAsRead(userID); err != nil {
		utils.ErrorLogger.Printf("Failed to mark all notifications as read: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update notifications"})
	}

	utils.InfoLogger.Printf("User %d marked all notifications as read", userID)
	return c.JSON(fiber.Map{"message": "All notifications marked as read"})
}

func (h *NotificationHandler) Delete(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	if err := h.service.Delete(uint(id), userID); err != nil {
		utils.ErrorLogger.Printf("Failed to delete notification: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete notification"})
	}

	return c.JSON(fiber.Map{"message": "Notification deleted successfully"})
}

func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	count, err := h.service.GetUnreadCount(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get unread count: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch unread count"})
	}

	return c.JSON(fiber.Map{"count": count})
}
