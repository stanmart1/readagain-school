package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type CartService struct {
	db *gorm.DB
}

func NewCartService(db *gorm.DB) *CartService {
	return &CartService{db: db}
}

type CartResponse struct {
	Items    []models.Cart `json:"items"`
	Subtotal float64       `json:"subtotal"`
	Total    float64       `json:"total"`
	Count    int           `json:"count"`
}

func (s *CartService) GetUserCart(userID uint) (*CartResponse, error) {
	var cartItems []models.Cart
	if err := s.db.Where("user_id = ?", userID).
		Preload("Book.Author.User").
		Preload("Book.Category").
		Find(&cartItems).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch cart", err)
	}

	response := &CartResponse{
		Items: cartItems,
		Count: len(cartItems),
	}

	for _, item := range cartItems {
		if item.Book != nil {
			response.Subtotal += item.Book.Price
		}
	}
	response.Total = response.Subtotal

	return response, nil
}

func (s *CartService) AddToCart(userID, bookID uint) (*CartResponse, error) {
	var book models.Book
	if err := s.db.First(&book, bookID).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found")
	}

	var userLibrary models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&userLibrary).Error; err == nil {
		return nil, utils.NewBadRequestError("You already own this book")
	}

	var existingCart models.Cart
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&existingCart).Error; err == nil {
		return nil, utils.NewBadRequestError("Book already in cart")
	}

	cart := models.Cart{
		UserID:   userID,
		BookID:   bookID,
		Quantity: 1,
	}

	if err := s.db.Create(&cart).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to add item to cart", err)
	}

	return s.GetUserCart(userID)
}

func (s *CartService) RemoveFromCart(userID, cartID uint) (*CartResponse, error) {
	var cart models.Cart
	if err := s.db.Where("id = ? AND user_id = ?", cartID, userID).First(&cart).Error; err != nil {
		return nil, utils.NewNotFoundError("Cart item not found")
	}

	if err := s.db.Delete(&cart).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to remove item from cart", err)
	}

	return s.GetUserCart(userID)
}

func (s *CartService) ClearCart(userID uint) error {
	if err := s.db.Where("user_id = ?", userID).Delete(&models.Cart{}).Error; err != nil {
		return utils.NewInternalServerError("Failed to clear cart", err)
	}
	return nil
}

func (s *CartService) GetCartCount(userID uint) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Cart{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, utils.NewInternalServerError("Failed to count cart items", err)
	}
	return count, nil
}

func (s *CartService) MergeGuestCart(userID uint, guestCartItems []GuestCartItem) (*CartResponse, error) {
	var ownedBookIDs []uint
	var userLibraries []models.UserLibrary
	if err := s.db.Where("user_id = ?", userID).Find(&userLibraries).Error; err == nil {
		for _, lib := range userLibraries {
			ownedBookIDs = append(ownedBookIDs, lib.BookID)
		}
	}

	var existingCartBookIDs []uint
	var existingCarts []models.Cart
	if err := s.db.Where("user_id = ?", userID).Find(&existingCarts).Error; err == nil {
		for _, cart := range existingCarts {
			existingCartBookIDs = append(existingCartBookIDs, cart.BookID)
		}
	}

	for _, guestItem := range guestCartItems {
		if contains(ownedBookIDs, guestItem.BookID) {
			continue
		}

		if contains(existingCartBookIDs, guestItem.BookID) {
			continue
		}

		var book models.Book
		if err := s.db.First(&book, guestItem.BookID).Error; err != nil {
			continue
		}

		cart := models.Cart{
			UserID:   userID,
			BookID:   guestItem.BookID,
			Quantity: 1,
		}

		s.db.Create(&cart)
	}

	return s.GetUserCart(userID)
}

func contains(slice []uint, item uint) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

type GuestCartItem struct {
	BookID uint `json:"book_id"`
}
