package middleware

import (
	"github.com/gofiber/fiber/v2"

	"readagain/internal/database"
	"readagain/internal/models"
	"readagain/internal/utils"
)

func AuthorOnly(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.NewUnauthorizedError("User not authenticated")
	}

	var author models.Author
	if err := database.DB.Where("user_id = ?", userID).First(&author).Error; err != nil {
		return utils.NewForbiddenError("Author access required")
	}

	c.Locals("author_id", author.ID)
	c.Locals("author", &author)

	return c.Next()
}

func AuthorOrAdmin(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return utils.NewUnauthorizedError("User not authenticated")
	}

	roleID, ok := c.Locals("role_id").(uint)
	if ok && (roleID == 1 || roleID == 2) {
		return c.Next()
	}

	var author models.Author
	if err := database.DB.Where("user_id = ?", userID).First(&author).Error; err != nil {
		return utils.NewForbiddenError("Author or admin access required")
	}

	c.Locals("author_id", author.ID)
	c.Locals("author", &author)

	return c.Next()
}
