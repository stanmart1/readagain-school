package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
)

type TestimonialService struct {
	db *gorm.DB
}

func NewTestimonialService(db *gorm.DB) *TestimonialService {
	return &TestimonialService{db: db}
}

func (s *TestimonialService) List(activeOnly bool) ([]models.Testimonial, error) {
	var testimonials []models.Testimonial
	query := s.db.Model(&models.Testimonial{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	if err := query.Order("order ASC, created_at DESC").Find(&testimonials).Error; err != nil {
		return nil, err
	}

	return testimonials, nil
}

func (s *TestimonialService) GetByID(id uint) (*models.Testimonial, error) {
	var testimonial models.Testimonial
	if err := s.db.First(&testimonial, id).Error; err != nil {
		return nil, err
	}
	return &testimonial, nil
}

func (s *TestimonialService) Create(testimonial *models.Testimonial) error {
	return s.db.Create(testimonial).Error
}

func (s *TestimonialService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Testimonial{}).Where("id = ?", id).Updates(updates).Error
}

func (s *TestimonialService) Delete(id uint) error {
	return s.db.Delete(&models.Testimonial{}, id).Error
}
