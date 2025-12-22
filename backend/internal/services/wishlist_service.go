package services

import (
	"errors"

	"gorm.io/gorm"

	"readagain/internal/models"
)

type WishlistService struct {
	db *gorm.DB
}

func NewWishlistService(db *gorm.DB) *WishlistService {
	return &WishlistService{db: db}
}

func (s *WishlistService) GetUserWishlist(userID uint) ([]models.Wishlist, error) {
	var wishlist []models.Wishlist
	if err := s.db.Where("user_id = ?", userID).Preload("Book").Order("created_at DESC").Find(&wishlist).Error; err != nil {
		return nil, err
	}
	return wishlist, nil
}

func (s *WishlistService) Add(userID, bookID uint) error {
	var existing models.Wishlist
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&existing).Error; err == nil {
		return errors.New("book already in wishlist")
	}

	var userBook models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&userBook).Error; err == nil {
		return errors.New("you already own this book")
	}

	wishlist := &models.Wishlist{
		UserID: userID,
		BookID: bookID,
	}
	return s.db.Create(wishlist).Error
}

func (s *WishlistService) Remove(userID, bookID uint) error {
	return s.db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&models.Wishlist{}).Error
}

func (s *WishlistService) RemoveByID(id, userID uint) error {
	return s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Wishlist{}).Error
}

func (s *WishlistService) IsInWishlist(userID, bookID uint) (bool, error) {
	var count int64
	if err := s.db.Model(&models.Wishlist{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
