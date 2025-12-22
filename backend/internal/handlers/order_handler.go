package handlers

import (
	"strconv"

	"readagain/internal/middleware"
	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) GetUserOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	status := c.Query("status", "")

	orders, meta, err := h.orderService.GetUserOrders(userID, page, limit, status)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get user orders: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve orders"})
	}

	return c.JSON(fiber.Map{
		"orders":     orders,
		"pagination": meta,
	})
}

func (h *OrderHandler) GetOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	order, err := h.orderService.GetOrderByID(uint(orderID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	if order.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	return c.JSON(fiber.Map{"order": order})
}

func (h *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	if err := h.orderService.CancelOrder(uint(orderID), userID); err != nil {
		utils.ErrorLogger.Printf("Failed to cancel order: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	middleware.LogAudit(c, "cancel_order", "order", uint(orderID), "", "")
	utils.InfoLogger.Printf("User %d cancelled order %d", userID, orderID)
	return c.JSON(fiber.Map{"message": "Order cancelled successfully"})
}

func (h *OrderHandler) GetAllOrders(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	status := c.Query("status", "")
	paymentMethod := c.Query("payment_method", "")

	orders, meta, err := h.orderService.GetAllOrders(page, limit, status, paymentMethod)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get all orders: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve orders"})
	}

	return c.JSON(fiber.Map{
		"orders":     orders,
		"pagination": meta,
	})
}

func (h *OrderHandler) GetOrderAdmin(c *fiber.Ctx) error {
	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	order, err := h.orderService.GetOrderByID(uint(orderID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	return c.JSON(fiber.Map{"order": order})
}

func (h *OrderHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order ID"})
	}

	var input struct {
		Status string `json:"status" validate:"required,oneof=pending processing completed cancelled failed"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	order, err := h.orderService.UpdateOrderStatus(uint(orderID), input.Status)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update order status: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Admin updated order %d status to %s", orderID, input.Status)
	return c.JSON(fiber.Map{"order": order})
}

func (h *OrderHandler) GetOrderStatistics(c *fiber.Ctx) error {
	stats, err := h.orderService.GetOrderStatistics()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get order statistics: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve statistics"})
	}

	return c.JSON(fiber.Map{"statistics": stats})
}
