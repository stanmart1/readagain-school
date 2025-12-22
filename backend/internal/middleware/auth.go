package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/config"
	"readagain/internal/database"
	"readagain/internal/models"
	"readagain/internal/utils"
)

func AuthRequired() fiber.Handler {
	cfg := config.Load()
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.NewUnauthorizedError("Missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return utils.NewUnauthorizedError("Invalid authorization format")
		}

		var blacklisted models.TokenBlacklist
		if err := database.DB.Where("token = ?", tokenString).First(&blacklisted).Error; err == nil {
			return utils.NewUnauthorizedError("Token has been revoked")
		}

		claims, err := utils.ValidateAccessToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			return utils.NewUnauthorizedError("Invalid or expired token")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("roleID", claims.RoleID)
		c.Locals("token", tokenString)

		return c.Next()
	}
}

func AdminRequired() fiber.Handler {
	cfg := config.Load()
	return func(c *fiber.Ctx) error {
		// Do auth check inline instead of calling AuthRequired()
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.NewUnauthorizedError("Missing authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return utils.NewUnauthorizedError("Invalid authorization format")
		}

		var blacklisted models.TokenBlacklist
		if err := database.DB.Where("token = ?", tokenString).First(&blacklisted).Error; err == nil {
			return utils.NewUnauthorizedError("Token has been revoked")
		}

		claims, err := utils.ValidateAccessToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			return utils.NewUnauthorizedError("Invalid or expired token")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("roleID", claims.RoleID)
		c.Locals("token", tokenString)

		// Check admin role
		roleID := claims.RoleID
		if roleID != 1 && roleID != 2 {
			return utils.NewForbiddenError("Admin access required")
		}

		return c.Next()
	}
}
