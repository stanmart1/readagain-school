package services

import (
	"gorm.io/gorm"
	"readagain/internal/models"
)

type WorkService struct {
	db *gorm.DB
}

func NewWorkService(db *gorm.DB) *WorkService {
	return &WorkService{db: db}
}

func (s *WorkService) GetAll() ([]models.Work, error) {
	var works []models.Work
	err := s.db.Where("is_active = ?", true).Order("\"order\" ASC, created_at DESC").Find(&works).Error
	return works, err
}

func (s *WorkService) GetAllAdmin() ([]models.Work, error) {
	var works []models.Work
	err := s.db.Order("\"order\" ASC, created_at DESC").Find(&works).Error
	return works, err
}

func (s *WorkService) Create(work *models.Work) error {
	return s.db.Create(work).Error
}

func (s *WorkService) Update(id uint, work *models.Work) error {
	return s.db.Model(&models.Work{}).Where("id = ?", id).Updates(work).Error
}

func (s *WorkService) Toggle(id uint) error {
	var work models.Work
	if err := s.db.First(&work, id).Error; err != nil {
		return err
	}
	return s.db.Model(&work).Update("is_active", !work.IsActive).Error
}

func (s *WorkService) Delete(id uint) error {
	return s.db.Delete(&models.Work{}, id).Error
}
