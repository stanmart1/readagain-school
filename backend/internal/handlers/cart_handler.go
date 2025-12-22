package handlers

import (
	"strconv"

	"readagain/internal/services"
	"readagain/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type CartHandler struct {
	cartService *services.CartService
}

func NewCartHandler(cartService *services.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	cart, err := h.cartService.GetUserCart(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get cart for user %d: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve cart"})
	}

	return c.JSON(fiber.Map{"cart": cart})
}

func (h *CartHandler) AddToCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		BookID uint `json:"book_id" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	cart, err := h.cartService.AddToCart(userID, input.BookID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to add to cart: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d added book %d to cart", userID, input.BookID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"cart": cart})
}

func (h *CartHandler) RemoveFromCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)
	cartItemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid cart item ID"})
	}

	cart, err := h.cartService.RemoveFromCart(userID, uint(cartItemID))
	if err != nil {
		utils.ErrorLogger.Printf("Failed to remove from cart: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	utils.InfoLogger.Printf("User %d removed item %d from cart", userID, cartItemID)
	return c.JSON(fiber.Map{"cart": cart})
}

func (h *CartHandler) ClearCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	if err := h.cartService.ClearCart(userID); err != nil {
		utils.ErrorLogger.Printf("Failed to clear cart: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to clear cart"})
	}

	utils.InfoLogger.Printf("User %d cleared cart", userID)
	return c.JSON(fiber.Map{"message": "Cart cleared successfully"})
}

func (h *CartHandler) GetCartCount(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	count, err := h.cartService.GetCartCount(userID)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to get cart count: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get cart count"})
	}

	return c.JSON(fiber.Map{"count": count})
}

func (h *CartHandler) MergeGuestCart(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint)

	var input struct {
		Items []services.GuestCartItem `json:"items" validate:"required"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := utils.Validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": utils.FormatValidationError(err)})
	}

	cart, err := h.cartService.MergeGuestCart(userID, input.Items)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to merge guest cart: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to merge cart"})
	}

	utils.InfoLogger.Printf("User %d merged guest cart with %d items", userID, len(input.Items))
	return c.JSON(fiber.Map{"cart": cart, "message": "Guest cart merged successfully"})
}
