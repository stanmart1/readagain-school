package handlers

import (
	"strconv"

	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) ListCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.ListCategories()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list categories: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve categories"})
	}

	return c.JSON(fiber.Map{"categories": categories})
}

func (h *CategoryHandler) GetCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	category, err := h.categoryService.GetCategoryByID(uint(categoryID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Category not found"})
	}

	return c.JSON(fiber.Map{"category": category})
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var input struct {
		Name        string `json:"name" validate:"required,min=2"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	category, err := h.categoryService.CreateCategory(input.Name, input.Description)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create category: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Created category: %s", input.Name)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"category": category})
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	category, err := h.categoryService.UpdateCategory(uint(categoryID), input.Name, input.Description)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to update category: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Updated category %d", categoryID)
	return c.JSON(fiber.Map{"category": category})
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid category ID"})
	}

	if err := h.categoryService.DeleteCategory(uint(categoryID)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete category: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("Deleted category %d", categoryID)
	return c.JSON(fiber.Map{"message": "Category deleted successfully"})
}
