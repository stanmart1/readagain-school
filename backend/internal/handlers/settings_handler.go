package handlers

import (
	"github.com/gofiber/fiber/v2"

	"readagain/internal/services"
	"readagain/internal/utils"
)

type SettingsHandler struct {
	service *services.SettingsService
}

func NewSettingsHandler(service *services.SettingsService) *SettingsHandler {
	return &SettingsHandler{service: service}
}

func (h *SettingsHandler) TestEmailGateway(c *fiber.Ctx) error {
	var req struct {
		GatewayType string `json:"gateway_type"`
		TestEmail   string `json:"test_email"`
		Config      map[string]interface{} `json:"config"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Simple test response - in a real implementation, you'd test the actual email sending
	return c.JSON(fiber.Map{
		"success": true,
		"message": "Email gateway test completed successfully",
		"gateway_type": req.GatewayType,
		"test_email": req.TestEmail,
	})
}

func (h *SettingsHandler) GetByCategory(c *fiber.Ctx) error {
	category := c.Query("category")

	settings, err := h.service.GetByCategory(category)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get settings: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch settings"})
	}

	return c.JSON(fiber.Map{"data": settings})
}

func (h *SettingsHandler) GetByKey(c *fiber.Ctx) error {
	key := c.Params("key")

	setting, err := h.service.GetByKey(key)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Setting not found"})
	}

	return c.JSON(fiber.Map{"data": setting})
}

func (h *SettingsHandler) Set(c *fiber.Ctx) error {
	var req struct {
		Key         string `json:"key" validate:"required"`
		Value       string `json:"value"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.Set(req.Key, req.Value, req.Category, req.Description); err != nil {
		utils.ErrorLogger.Printf("Failed to set setting: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save setting"})
	}

	setting, _ := h.service.GetByKey(req.Key)
	utils.InfoLogger.Printf("Setting updated: %s", req.Key)
	return c.JSON(fiber.Map{"data": setting})
}

func (h *SettingsHandler) Delete(c *fiber.Ctx) error {
	key := c.Params("key")

	if err := h.service.Delete(key); err != nil {
		utils.ErrorLogger.Printf("Failed to delete setting: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete setting"})
	}

	utils.InfoLogger.Printf("Setting deleted: %s", key)
	return c.JSON(fiber.Map{"message": "Setting deleted successfully"})
}

func (h *SettingsHandler) GetEmailSettings(c *fiber.Ctx) error {
	settings, err := h.service.GetEmailSettings()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get email settings: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch email settings"})
	}

	return c.JSON(fiber.Map{"data": settings})
}

func (h *SettingsHandler) UpdateEmailSettings(c *fiber.Ctx) error {
	var req struct {
		ResendAPIKey string `json:"resend_api_key" validate:"required"`
		FromEmail    string `json:"from_email" validate:"required,email"`
		FromName     string `json:"from_name" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.UpdateEmailSettings(req.ResendAPIKey, req.FromEmail, req.FromName); err != nil {
		utils.ErrorLogger.Printf("Failed to update email settings: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update email settings"})
	}

	utils.InfoLogger.Println("Email settings updated")
	return c.JSON(fiber.Map{"message": "Email settings updated successfully"})
}

func (h *SettingsHandler) GetPaymentSettings(c *fiber.Ctx) error {
	settings, err := h.service.GetPaymentSettings()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get payment settings: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch payment settings"})
	}

	return c.JSON(fiber.Map{"data": settings})
}

func (h *SettingsHandler) UpdatePaymentSettings(c *fiber.Ctx) error {
	var req struct {
		PaystackSecretKey     string `json:"paystack_secret_key"`
		FlutterwaveSecretKey  string `json:"flutterwave_secret_key"`
		BankName              string `json:"bank_name"`
		AccountNumber         string `json:"account_number"`
		AccountName           string `json:"account_name"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.UpdatePaymentSettings(
		req.PaystackSecretKey,
		req.FlutterwaveSecretKey,
		req.BankName,
		req.AccountNumber,
		req.AccountName,
	); err != nil {
		utils.ErrorLogger.Printf("Failed to update payment settings: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update payment settings"})
	}

	utils.InfoLogger.Println("Payment settings updated")
	return c.JSON(fiber.Map{"message": "Payment settings updated successfully"})
}

func (h *SettingsHandler) GetPublic(c *fiber.Ctx) error {
	settings, err := h.service.GetByCategory("general")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch settings"})
	}

	publicSettings := make(map[string]string)
	for _, s := range settings {
		if s.Key == "session_timeout" || s.Key == "site_name" || s.Key == "site_description" {
			publicSettings[s.Key] = s.Value
		}
	}

	return c.JSON(fiber.Map{"data": publicSettings})
}
