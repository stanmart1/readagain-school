package services

import (
	"math"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type AuditService struct {
	db *gorm.DB
}

func NewAuditService(db *gorm.DB) *AuditService {
	return &AuditService{db: db}
}

func (s *AuditService) Log(userID uint, action, entityType string, entityID uint, oldValue, newValue, ipAddress, userAgent string) error {
	audit := &models.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValue:   oldValue,
		NewValue:   newValue,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}
	return s.db.Create(audit).Error
}

func (s *AuditService) GetLogs(page, limit int, userID *uint, entityType, action string) ([]models.AuditLog, *utils.PaginationMeta, error) {
	var logs []models.AuditLog
	var total int64

	query := s.db.Model(&models.AuditLog{}).Preload("User")

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	if action != "" {
		query = query.Where("action = ?", action)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return logs, meta, nil
}

func (s *AuditService) GetByID(id uint) (*models.AuditLog, error) {
	var log models.AuditLog
	if err := s.db.Preload("User").First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}
