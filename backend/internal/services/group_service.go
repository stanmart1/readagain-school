package services

import (
	"errors"
	"readagain/internal/models"

	"gorm.io/gorm"
)

type GroupService struct {
	db *gorm.DB
}

func NewGroupService(db *gorm.DB) *GroupService {
	return &GroupService{db: db}
}

func (s *GroupService) GetAll() ([]models.Group, error) {
	var groups []models.Group
	err := s.db.Preload("Creator").Order("created_at DESC").Find(&groups).Error
	return groups, err
}

func (s *GroupService) GetByID(id uint) (*models.Group, error) {
	var group models.Group
	err := s.db.Preload("Creator").First(&group, id).Error
	return &group, err
}

func (s *GroupService) Create(group *models.Group) error {
	return s.db.Create(group).Error
}

func (s *GroupService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Group{}).Where("id = ?", id).Updates(updates).Error
}

func (s *GroupService) Delete(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ?", id).Delete(&models.GroupMember{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.Group{}, id).Error
	})
}

func (s *GroupService) GetMembers(groupID uint) ([]models.User, error) {
	var users []models.User
	err := s.db.Joins("JOIN group_members ON group_members.user_id = users.id").
		Where("group_members.group_id = ?", groupID).
		Preload("Role").
		Find(&users).Error
	return users, err
}

func (s *GroupService) AddMember(groupID, userID uint) error {
	var exists int64
	s.db.Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&exists)
	if exists > 0 {
		return errors.New("user already in group")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		member := &models.GroupMember{GroupID: groupID, UserID: userID}
		if err := tx.Create(member).Error; err != nil {
			return err
		}
		return tx.Model(&models.Group{}).Where("id = ?", groupID).UpdateColumn("member_count", gorm.Expr("member_count + ?", 1)).Error
	})
}

func (s *GroupService) RemoveMember(groupID, userID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("group_id = ? AND user_id = ?", groupID, userID).Delete(&models.GroupMember{}).Error; err != nil {
			return err
		}
		return tx.Model(&models.Group{}).Where("id = ?", groupID).UpdateColumn("member_count", gorm.Expr("member_count - ?", 1)).Error
	})
}

func (s *GroupService) AddMembers(groupID uint, userIDs []uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, userID := range userIDs {
			var exists int64
			tx.Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Count(&exists)
			if exists == 0 {
				member := &models.GroupMember{GroupID: groupID, UserID: userID}
				if err := tx.Create(member).Error; err != nil {
					return err
				}
			}
		}
		var count int64
		tx.Model(&models.GroupMember{}).Where("group_id = ?", groupID).Count(&count)
		return tx.Model(&models.Group{}).Where("id = ?", groupID).Update("member_count", count).Error
	})
}
