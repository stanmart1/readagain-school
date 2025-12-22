package handlers

import (
	"fmt"

	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CheckoutHandler struct {
	orderService   *services.OrderService
	paymentService *services.PaymentService
}

func NewCheckoutHandler(orderService *services.OrderService, paymentService *services.PaymentService) *CheckoutHandler {
	return &CheckoutHandler{
		orderService:   orderService,
		paymentService: paymentService,
	}
}

func (h *CheckoutHandler) InitializeCheckout(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		PaymentMethod string `json:"payment_method" validate:"required,oneof=paystack flutterwave bank_transfer"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	order, err := h.orderService.CreateOrderFromCart(userID, input.PaymentMethod)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create order: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d created order %s", userID, order.OrderNumber)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"order": order})
}

func (h *CheckoutHandler) InitializePayment(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	email := c.Locals("email").(string)

	var input struct {
		OrderID uint   `json:"order_id" validate:"required"`
		Name    string `json:"name"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	order, err := h.orderService.GetOrderByID(input.OrderID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	if order.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied"})
	}

	if order.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order is not pending"})
	}

	paymentReq := services.PaymentInitializeRequest{
		Provider:  order.PaymentMethod,
		Email:     email,
		Name:      input.Name,
		Amount:    order.TotalAmount,
		Reference: order.Notes,
		Metadata: map[string]interface{}{
			"order_id":     order.ID,
			"order_number": order.OrderNumber,
			"user_id":      userID,
		},
	}

	paymentResp, err := h.paymentService.InitializePayment(paymentReq)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to initialize payment: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Payment initialized for order %s", order.OrderNumber)
	return c.JSON(fiber.Map{"payment": paymentResp})
}

func (h *CheckoutHandler) VerifyPayment(c *fiber.Ctx) error {
	reference := c.Params("reference")
	provider := c.Query("provider", "paystack")

	paymentResp, err := h.paymentService.VerifyPayment(provider, reference)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to verify payment: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if paymentResp.Status == "success" || paymentResp.Status == "successful" {
		// Payment verified, order will be completed via webhook
		utils.InfoLogger.Printf("Payment verified: %s", reference)
	}

	return c.JSON(fiber.Map{"payment": paymentResp})
}

func (h *CheckoutHandler) PaystackWebhook(c *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	event, ok := payload["event"].(string)
	if !ok || event != "charge.success" {
		return c.SendStatus(fiber.StatusOK)
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		return c.SendStatus(fiber.StatusOK)
	}

	reference, _ := data["reference"].(string)
	status, _ := data["status"].(string)

	if status == "success" {
		order, err := h.orderService.GetOrderByReference(reference)
		if err == nil {
			transactionID := ""
			if id, ok := data["id"].(float64); ok {
				transactionID = fmt.Sprintf("%.0f", id)
			}
			h.orderService.CompleteOrder(order.ID, transactionID)
		}
	}

	utils.InfoLogger.Printf("Paystack webhook processed: %s - %s", reference, status)
	return c.SendStatus(fiber.StatusOK)
}

func (h *CheckoutHandler) FlutterwaveWebhook(c *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid payload"})
	}

	event, ok := payload["event"].(string)
	if !ok || event != "charge.completed" {
		return c.SendStatus(fiber.StatusOK)
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		return c.SendStatus(fiber.StatusOK)
	}

	txRef, _ := data["tx_ref"].(string)
	status, _ := data["status"].(string)

	if status == "successful" {
		order, err := h.orderService.GetOrderByReference(txRef)
		if err == nil {
			transactionID := ""
			if id, ok := data["id"].(float64); ok {
				transactionID = fmt.Sprintf("%.0f", id)
			}
			h.orderService.CompleteOrder(order.ID, transactionID)
		}
	}

	utils.InfoLogger.Printf("Flutterwave webhook processed: %s - %s", txRef, status)
	return c.SendStatus(fiber.StatusOK)
}
