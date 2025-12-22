package handlers

import (
	"strconv"

	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthorHandler struct {
	authorService *services.AuthorService
}

func NewAuthorHandler(authorService *services.AuthorService) *AuthorHandler {
	return &AuthorHandler{authorService: authorService}
}

func (h *AuthorHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.authorService.GetStats()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get author stats: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch author statistics"})
	}

	return c.JSON(stats)
}

func (h *AuthorHandler) ListAuthors(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	search := c.Query("search", "")

	authors, meta, err := h.authorService.ListAuthors(page, limit, search)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list authors: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve authors"})
	}

	return c.JSON(fiber.Map{
		"authors":    authors,
		"pagination": meta,
	})
}

func (h *AuthorHandler) GetAuthor(c *fiber.Ctx) error {
	authorID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid author ID"})
	}

	author, err := h.authorService.GetAuthorByID(uint(authorID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Author not found"})
	}

	return c.JSON(fiber.Map{"author": author})
}

func (h *AuthorHandler) CreateAuthor(c *fiber.Ctx) error {
	var input struct {
		UserID  uint   `json:"user_id" validate:"required"`
		PenName string `json:"pen_name"`
		Bio     string `json:"bio"`
		Website string `json:"website"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	author, err := h.authorService.CreateAuthor(input.UserID, input.PenName, input.Bio, input.Website)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create author: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Created author for user %d", input.UserID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"author": author})
}

func (h *AuthorHandler) UpdateAuthor(c *fiber.Ctx) error {
	authorID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid author ID"})
	}

	var input map[string]interface{}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	author, err := h.authorService.UpdateAuthor(uint(authorID), input)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update author: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Updated author %d", authorID)
	return c.JSON(fiber.Map{"author": author})
}

func (h *AuthorHandler) DeleteAuthor(c *fiber.Ctx) error {
	authorID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid author ID"})
	}

	if err := h.authorService.DeleteAuthor(uint(authorID)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete author: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Deleted author %d", authorID)
	return c.JSON(fiber.Map{"message": "Author deleted successfully"})
}
