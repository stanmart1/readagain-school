package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/middleware"
	"readagain/internal/models"
	"readagain/internal/services"
	"readagain/internal/utils"
)

type ReviewHandler struct {
	service *services.ReviewService
}

func NewReviewHandler(service *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	var review models.Review
	if err := c.BodyParser(&review); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&review); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	review.UserID = userID
	review.BookID = uint(bookID)

	if err := h.service.Create(&review); err != nil {
		utils.ErrorLogger.Printf("Failed to create review: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d created review for book %d", userID, bookID)
	return c.Status(201).JSON(fiber.Map{"data": review})
}

func (h *ReviewHandler) GetBookReviews(c *fiber.Ctx) error {
	bookID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	reviews, meta, err := h.service.GetBookReviews(uint(bookID), page, limit, true)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get book reviews: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch reviews"})
	}

	avgRating, count, _ := h.service.GetBookRating(uint(bookID))

	return c.JSON(fiber.Map{
		"data":          reviews,
		"meta":          meta,
		"average_rating": avgRating,
		"review_count":   count,
	})
}

func (h *ReviewHandler) ListAllReviews(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	status := c.Query("status")

	reviews, meta, err := h.service.ListAll(page, limit, status)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to list reviews: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch reviews"})
	}

	return c.JSON(fiber.Map{"reviews": reviews, "meta": meta})
}

func (h *ReviewHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.service.GetStats()
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get review stats: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch stats"})
	}

	return c.JSON(fiber.Map{"stats": stats})
}

func (h *ReviewHandler) UpdateStatus(c *fiber.Ctx) error {
	var req struct {
		ReviewID uint   `json:"reviewId" validate:"required"`
		Status   string `json:"status" validate:"required,oneof=pending approved rejected"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.UpdateStatus(req.ReviewID, req.Status); err != nil {
		utils.ErrorLogger.Printf("Failed to update review status: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update review"})
	}

	middleware.LogAudit(c, "update_review_status", "review", req.ReviewID, "", req.Status)
	utils.InfoLogger.Printf("Review %d status updated to %s", req.ReviewID, req.Status)
	return c.JSON(fiber.Map{"message": "Review status updated"})
}

func (h *ReviewHandler) ToggleFeatured(c *fiber.Ctx) error {
	var req struct {
		ReviewID   uint `json:"reviewId" validate:"required"`
		IsFeatured bool `json:"isFeatured"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.ToggleFeatured(req.ReviewID, req.IsFeatured); err != nil {
		utils.ErrorLogger.Printf("Failed to toggle featured: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update review"})
	}

	middleware.LogAudit(c, "toggle_review_featured", "review", req.ReviewID, "", "")
	utils.InfoLogger.Printf("Review %d featured status: %v", req.ReviewID, req.IsFeatured)
	return c.JSON(fiber.Map{"message": "Review featured status updated"})
}

func (h *ReviewHandler) Delete(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid review ID"})
	}

	if err := h.service.Delete(uint(id)); err != nil {
		utils.ErrorLogger.Printf("Failed to delete review: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete review"})
	}

	middleware.LogAudit(c, "delete_review", "review", uint(id), "", "")
	utils.InfoLogger.Printf("Review deleted: %d", id)
	return c.JSON(fiber.Map{"message": "Review deleted successfully"})
}

func (h *ReviewHandler) GetFeatured(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	reviews, err := h.service.GetFeatured(limit)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get featured reviews: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch featured reviews"})
	}

	return c.JSON(fiber.Map{"data": reviews})
}
