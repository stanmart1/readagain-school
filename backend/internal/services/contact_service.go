package services

import (
	"math"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type ContactService struct {
	db *gorm.DB
}

func NewContactService(db *gorm.DB) *ContactService {
	return &ContactService{db: db}
}

func (s *ContactService) Create(message *models.ContactMessage) error {
	return s.db.Create(message).Error
}

func (s *ContactService) List(page, limit int, status string) ([]models.ContactMessage, *utils.PaginationMeta, error) {
	var messages []models.ContactMessage
	var total int64

	query := s.db.Model(&models.ContactMessage{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return messages, meta, nil
}

func (s *ContactService) GetByID(id uint) (*models.ContactMessage, error) {
	var message models.ContactMessage
	if err := s.db.First(&message, id).Error; err != nil {
		return nil, err
	}
	return &message, nil
}

func (s *ContactService) Reply(id uint, reply string) error {
	return s.db.Model(&models.ContactMessage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"reply":  reply,
		"status": "replied",
	}).Error
}

func (s *ContactService) UpdateStatus(id uint, status string) error {
	return s.db.Model(&models.ContactMessage{}).Where("id = ?", id).Update("status", status).Error
}

func (s *ContactService) Delete(id uint) error {
	return s.db.Delete(&models.ContactMessage{}, id).Error
}
