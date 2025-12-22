package handlers

import (
	"github.com/gofiber/fiber/v2"

	"readagain/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" validate:"required"`
	Password        string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := h.authService.Register(req.Email, req.Username, req.Password, req.FirstName, req.LastName)
	if err != nil {
		return err
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User registered successfully",
		"user":    user,
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	accessToken, refreshToken, user, err := h.authService.Login(req.EmailOrUsername, req.Password, ipAddress, userAgent)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	accessToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := c.Locals("token").(string)
	userID := c.Locals("userID").(uint)

	if err := h.authService.Logout(token, userID); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	email := c.Locals("email").(string)
	roleID := c.Locals("roleID").(uint)

	return c.JSON(fiber.Map{
		"user_id": userID,
		"email":   email,
		"role_id": roleID,
	})
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	resetToken, err := h.authService.RequestPasswordReset(req.Email)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Password reset token sent",
		"token":   resetToken,
	})
}

func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req ResetPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Password reset successfully",
	})
}
