package services

import (
	"errors"
	"math"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type ReviewService struct {
	db *gorm.DB
}

func NewReviewService(db *gorm.DB) *ReviewService {
	return &ReviewService{db: db}
}

func (s *ReviewService) Create(review *models.Review) error {
	var existingReview models.Review
	err := s.db.Where("user_id = ? AND book_id = ?", review.UserID, review.BookID).First(&existingReview).Error
	if err == nil {
		return errors.New("you have already reviewed this book")
	}

	var userBook models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", review.UserID, review.BookID).First(&userBook).Error; err != nil {
		return errors.New("you can only review books you own")
	}

	review.Status = "pending"
	return s.db.Create(review).Error
}

func (s *ReviewService) GetBookReviews(bookID uint, page, limit int, approvedOnly bool) ([]models.Review, *utils.PaginationMeta, error) {
	var reviews []models.Review
	var total int64

	query := s.db.Model(&models.Review{}).Where("book_id = ?", bookID)

	if approvedOnly {
		query = query.Where("status = ?", "approved")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return reviews, meta, nil
}

func (s *ReviewService) ListAll(page, limit int, status string) ([]models.Review, *utils.PaginationMeta, error) {
	var reviews []models.Review
	var total int64

	query := s.db.Model(&models.Review{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("User").Preload("Book").Order("created_at DESC").Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return reviews, meta, nil
}

func (s *ReviewService) GetStats() (map[string]interface{}, error) {
	var total, pending, approved, rejected, featured int64
	var avgRating float64

	s.db.Model(&models.Review{}).Count(&total)
	s.db.Model(&models.Review{}).Where("status = ?", "pending").Count(&pending)
	s.db.Model(&models.Review{}).Where("status = ?", "approved").Count(&approved)
	s.db.Model(&models.Review{}).Where("status = ?", "rejected").Count(&rejected)
	s.db.Model(&models.Review{}).Where("is_featured = ?", true).Count(&featured)
	s.db.Model(&models.Review{}).Where("status = ?", "approved").Select("COALESCE(AVG(rating), 0)").Scan(&avgRating)

	return map[string]interface{}{
		"total":         total,
		"pending":       pending,
		"approved":      approved,
		"rejected":      rejected,
		"featured":      featured,
		"averageRating": avgRating,
	}, nil
}

func (s *ReviewService) UpdateStatus(id uint, status string) error {
	return s.db.Model(&models.Review{}).Where("id = ?", id).Update("status", status).Error
}

func (s *ReviewService) ToggleFeatured(id uint, isFeatured bool) error {
	return s.db.Model(&models.Review{}).Where("id = ?", id).Update("is_featured", isFeatured).Error
}

func (s *ReviewService) Delete(id uint) error {
	return s.db.Delete(&models.Review{}, id).Error
}

func (s *ReviewService) GetBookRating(bookID uint) (float64, int64, error) {
	var avgRating float64
	var count int64

	s.db.Model(&models.Review{}).Where("book_id = ? AND status = ?", bookID, "approved").Count(&count)
	s.db.Model(&models.Review{}).Where("book_id = ? AND status = ?", bookID, "approved").Select("AVG(rating)").Scan(&avgRating)

	return avgRating, count, nil
}

func (s *ReviewService) GetFeatured(limit int) ([]models.Review, error) {
	var reviews []models.Review
	err := s.db.Where("is_featured = ? AND status = ?", true, "approved").
		Preload("User").
		Preload("Book").
		Order("created_at DESC").
		Limit(limit).
		Find(&reviews).Error
	return reviews, err
}
