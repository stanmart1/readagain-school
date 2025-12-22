package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"readagain/internal/services"
	"readagain/internal/utils"
)

type WishlistHandler struct {
	service *services.WishlistService
}

func NewWishlistHandler(service *services.WishlistService) *WishlistHandler {
	return &WishlistHandler{service: service}
}

func (h *WishlistHandler) GetWishlist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	wishlist, err := h.service.GetUserWishlist(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get wishlist: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch wishlist"})
	}

	return c.JSON(fiber.Map{"data": wishlist})
}

func (h *WishlistHandler) AddToWishlist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var req struct {
		BookID uint `json:"book_id" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := utils.Validate.Struct(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	if err := h.service.Add(userID, req.BookID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d added book %d to wishlist", userID, req.BookID)
	return c.Status(201).JSON(fiber.Map{"message": "Book added to wishlist"})
}

func (h *WishlistHandler) RemoveFromWishlist(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid wishlist ID"})
	}

	if err := h.service.RemoveByID(uint(id), userID); err != nil {
		utils.ErrorLogger.Printf("Failed to remove from wishlist: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to remove from wishlist"})
	}

	utils.InfoLogger.Printf("User %d removed item %d from wishlist", userID, id)
	return c.JSON(fiber.Map{"message": "Book removed from wishlist"})
}
