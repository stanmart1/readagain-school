package middleware

import (
	"github.com/gofiber/fiber/v2"

	"readagain/internal/database"
	"readagain/internal/models"
	"readagain/internal/utils"
)

func RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := AuthRequired()(c); err != nil {
			return err
		}

		userID := c.Locals("userID").(uint)

		var user models.User
		if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
			return utils.NewUnauthorizedError("User not found")
		}

		if user.Role == nil {
			return utils.NewForbiddenError("No role assigned")
		}

		hasPermission := false
		for _, rp := range user.Role.Permissions {
			if rp.Name == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			return utils.NewForbiddenError("Insufficient permissions")
		}

		return c.Next()
	}
}

func RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := AuthRequired()(c); err != nil {
			return err
		}

		userID := c.Locals("userID").(uint)

		var user models.User
		if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
			return utils.NewUnauthorizedError("User not found")
		}

		if user.Role == nil {
			return utils.NewForbiddenError("No role assigned")
		}

		for _, rp := range user.Role.Permissions {
			for _, perm := range permissions {
				if rp.Name == perm {
					return c.Next()
				}
			}
		}

		return utils.NewForbiddenError("Insufficient permissions")
	}
}

func RequireAllPermissions(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := AuthRequired()(c); err != nil {
			return err
		}

		userID := c.Locals("userID").(uint)

		var user models.User
		if err := database.DB.Preload("Role.Permissions").First(&user, userID).Error; err != nil {
			return utils.NewUnauthorizedError("User not found")
		}

		if user.Role == nil {
			return utils.NewForbiddenError("No role assigned")
		}

		userPermissions := make(map[string]bool)
		for _, rp := range user.Role.Permissions {
			userPermissions[rp.Name] = true
		}

		for _, perm := range permissions {
			if !userPermissions[perm] {
				return utils.NewForbiddenError("Insufficient permissions")
			}
		}

		return c.Next()
	}
}

func RequireRole(roleName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := AuthRequired()(c); err != nil {
			return err
		}

		userID := c.Locals("userID").(uint)

		var user models.User
		if err := database.DB.Preload("Role").First(&user, userID).Error; err != nil {
			return utils.NewUnauthorizedError("User not found")
		}

		if user.Role == nil || user.Role.Name != roleName {
			return utils.NewForbiddenError("Role required: " + roleName)
		}

		return c.Next()
	}
}
