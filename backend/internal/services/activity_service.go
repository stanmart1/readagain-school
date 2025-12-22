package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
)

type ActivityService struct {
	db *gorm.DB
}

func NewActivityService(db *gorm.DB) *ActivityService {
	return &ActivityService{db: db}
}

func (s *ActivityService) Log(userID uint, activityType, title, description, entityType string, entityID uint, metadata string) error {
	activity := &models.Activity{
		UserID:      userID,
		Type:        activityType,
		Title:       title,
		Description: description,
		EntityType:  entityType,
		EntityID:    entityID,
		Metadata:    metadata,
	}
	return s.db.Create(activity).Error
}

func (s *ActivityService) GetUserActivities(userID uint, limit, offset int) ([]models.Activity, bool, error) {
	var activities []models.Activity

	if err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit + 1).
		Offset(offset).
		Find(&activities).Error; err != nil {
		return nil, false, err
	}

	hasMore := len(activities) > limit
	if hasMore {
		activities = activities[:limit]
	}

	return activities, hasMore, nil
}
