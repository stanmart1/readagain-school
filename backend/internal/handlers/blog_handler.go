package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/models"
	"readagain/internal/services"
	"readagain/internal/utils"
)

type BlogHandler struct {
	service *services.BlogService
}

func NewBlogHandler(service *services.BlogService) *BlogHandler {
	return &BlogHandler{service: service}
}

func (h *BlogHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.service.GetStats()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get blog stats: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch blog statistics"})
	}

	return c.JSON(stats)
}

func (h *BlogHandler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status", "published")
	search := c.Query("search")
	tag := c.Query("tag")

	blogs, meta, err := h.service.List(page, limit, status, search, tag)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list blogs: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch blogs"})
	}

	return c.JSON(fiber.Map{"data": blogs, "meta": meta})
}

func (h *BlogHandler) AdminList(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	status := c.Query("status")
	search := c.Query("search")
	tag := c.Query("tag")

	blogs, meta, err := h.service.List(page, limit, status, search, tag)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list blogs: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch blogs"})
	}

	return c.JSON(fiber.Map{"data": blogs, "meta": meta})
}

func (h *BlogHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	blog, err := h.service.GetBySlug(slug)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Blog not found"})
	}

	if blog.Status == "published" {
		go h.service.IncrementViews(blog.ID)
	}

	return c.JSON(fiber.Map{"data": blog})
}

func (h *BlogHandler) GetByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid blog ID"})
	}

	blog, err := h.service.GetByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Blog not found"})
	}

	return c.JSON(fiber.Map{"data": blog})
}

func (h *BlogHandler) Create(c *fiber.Ctx) error {
	var blog models.Blog
	if err := c.BodyParser(&blog); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&blog); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	userID := c.Locals("userID").(uint)
	blog.AuthorID = userID

	if err := h.service.Create(&blog); err != nil {
		utils.ErrorLogger.Printf("Failed to create blog: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create blog"})
	}

	utils.InfoLogger.Printf("Blog created: %s by user %d", blog.Title, userID)
	return c.Status(201).JSON(fiber.Map{"data": blog})
}

func (h *BlogHandler) Update(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid blog ID"})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.Update(uint(id), updates); err != nil {
		utils.ErrorLogger.Printf("Failed to update blog: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update blog"})
	}

	blog, _ := h.service.GetByID(uint(id))
	utils.InfoLogger.Printf("Blog updated: %d", id)
	return c.JSON(fiber.Map{"data": blog})
}

func (h *BlogHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid blog ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete blog: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete blog"})
	}

	utils.InfoLogger.Printf("Blog deleted: %d", id)
	return c.JSON(fiber.Map{"message": "Blog deleted successfully"})
}
