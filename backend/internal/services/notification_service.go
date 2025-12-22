package services

import (
	"math"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

func (s *NotificationService) Create(notification *models.Notification) error {
	return s.db.Create(notification).Error
}

func (s *NotificationService) GetUserNotifications(userID uint, page, limit int, unreadOnly bool) ([]models.Notification, *utils.PaginationMeta, error) {
	var notifications []models.Notification
	var total int64

	query := s.db.Model(&models.Notification{}).Where("user_id = ?", userID)

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * limit
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&notifications).Error; err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return notifications, meta, nil
}

func (s *NotificationService) GetByID(id, userID uint) (*models.Notification, error) {
	var notification models.Notification
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&notification).Error; err != nil {
		return nil, err
	}
	return &notification, nil
}

func (s *NotificationService) MarkAsRead(id, userID uint) error {
	return s.db.Model(&models.Notification{}).Where("id = ? AND user_id = ?", id, userID).Update("is_read", true).Error
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true).Error
}

func (s *NotificationService) Delete(id, userID uint) error {
	return s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Notification{}).Error
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	if err := s.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (s *NotificationService) Notify(userID uint, notifType, title, message, actionURL string) error {
	notification := &models.Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Message:   message,
		ActionURL: actionURL,
	}
	return s.Create(notification)
}
