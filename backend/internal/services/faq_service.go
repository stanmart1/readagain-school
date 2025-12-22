package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
)

type FAQService struct {
	db *gorm.DB
}

func NewFAQService(db *gorm.DB) *FAQService {
	return &FAQService{db: db}
}

func (s *FAQService) List(category string, activeOnly bool) ([]models.FAQ, error) {
	var faqs []models.FAQ
	query := s.db.Model(&models.FAQ{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Order("\"order\" ASC, created_at ASC").Find(&faqs).Error; err != nil {
		return nil, err
	}

	return faqs, nil
}

func (s *FAQService) GetByID(id uint) (*models.FAQ, error) {
	var faq models.FAQ
	if err := s.db.First(&faq, id).Error; err != nil {
		return nil, err
	}
	return &faq, nil
}

func (s *FAQService) Create(faq *models.FAQ) error {
	return s.db.Create(faq).Error
}

func (s *FAQService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.FAQ{}).Where("id = ?", id).Updates(updates).Error
}

func (s *FAQService) Delete(id uint) error {
	return s.db.Delete(&models.FAQ{}, id).Error
}

func (s *FAQService) GetCategories() ([]string, error) {
	var categories []string
	if err := s.db.Model(&models.FAQ{}).Distinct("category").Where("category IS NOT NULL AND category != ''").Pluck("category", &categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
