package services

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"readagain/internal/models"
)

type ChatService struct {
	db *gorm.DB
}

func NewChatService(db *gorm.DB) *ChatService {
	return &ChatService{db: db}
}

// Room Management
func (s *ChatService) CreateRoom(room *models.ChatRoom) error {
	return s.db.Create(room).Error
}

func (s *ChatService) GetRoomByID(roomID uint) (*models.ChatRoom, error) {
	var room models.ChatRoom
	err := s.db.Preload("Group").Preload("Book").Preload("Creator").First(&room, roomID).Error
	return &room, err
}

func (s *ChatService) GetRoomsByType(roomType string, page, limit int) ([]models.ChatRoom, int64, error) {
	var rooms []models.ChatRoom
	var total int64

	query := s.db.Model(&models.ChatRoom{}).Where("type = ? AND is_active = ?", roomType, true)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Group").Preload("Book").Preload("Creator").
		Order("last_message_at DESC NULLS LAST, created_at DESC").
		Offset(offset).Limit(limit).Find(&rooms).Error

	return rooms, total, err
}

func (s *ChatService) GetUserRooms(userID uint, page, limit int) ([]models.ChatRoom, int64, error) {
	var rooms []models.ChatRoom
	var total int64

	subQuery := s.db.Model(&models.ChatMember{}).Select("room_id").Where("user_id = ?", userID)
	
	query := s.db.Model(&models.ChatRoom{}).Where("id IN (?) AND is_active = ?", subQuery, true)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("Group").Preload("Book").Preload("Creator").
		Order("last_message_at DESC NULLS LAST, created_at DESC").
		Offset(offset).Limit(limit).Find(&rooms).Error

	return rooms, total, err
}

func (s *ChatService) UpdateRoom(roomID uint, updates map[string]interface{}) error {
	return s.db.Model(&models.ChatRoom{}).Where("id = ?", roomID).Updates(updates).Error
}

func (s *ChatService) DeleteRoom(roomID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete all messages
		if err := tx.Where("room_id = ?", roomID).Delete(&models.ChatMessage{}).Error; err != nil {
			return err
		}
		// Delete all members
		if err := tx.Where("room_id = ?", roomID).Delete(&models.ChatMember{}).Error; err != nil {
			return err
		}
		// Delete room
		return tx.Delete(&models.ChatRoom{}, roomID).Error
	})
}

// Member Management
func (s *ChatService) AddMember(member *models.ChatMember) error {
	member.JoinedAt = time.Now()
	return s.db.Create(member).Error
}

func (s *ChatService) AddMembers(members []models.ChatMember) error {
	now := time.Now()
	for i := range members {
		members[i].JoinedAt = now
	}
	return s.db.Create(&members).Error
}

func (s *ChatService) RemoveMember(roomID, userID uint) error {
	return s.db.Where("room_id = ? AND user_id = ?", roomID, userID).Delete(&models.ChatMember{}).Error
}

func (s *ChatService) GetRoomMembers(roomID uint) ([]models.ChatMember, error) {
	var members []models.ChatMember
	err := s.db.Where("room_id = ?", roomID).Preload("User").Find(&members).Error
	return members, err
}

func (s *ChatService) IsMember(roomID, userID uint) (bool, error) {
	var count int64
	err := s.db.Model(&models.ChatMember{}).Where("room_id = ? AND user_id = ?", roomID, userID).Count(&count).Error
	return count > 0, err
}

func (s *ChatService) UpdateMemberRole(roomID, userID uint, role string) error {
	return s.db.Model(&models.ChatMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Update("role", role).Error
}

func (s *ChatService) UpdateLastRead(roomID, userID uint) error {
	now := time.Now()
	return s.db.Model(&models.ChatMember{}).
		Where("room_id = ? AND user_id = ?", roomID, userID).
		Updates(map[string]interface{}{
			"last_read_at":  &now,
			"unread_count": 0,
		}).Error
}

// Message Management
func (s *ChatService) CreateMessage(message *models.ChatMessage) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Create message
		if err := tx.Create(message).Error; err != nil {
			return err
		}

		// Update room's last message
		now := time.Now()
		if err := tx.Model(&models.ChatRoom{}).Where("id = ?", message.RoomID).Updates(map[string]interface{}{
			"last_message":    message.Message,
			"last_message_at": &now,
		}).Error; err != nil {
			return err
		}

		// Increment unread count for all members except sender
		return tx.Model(&models.ChatMember{}).
			Where("room_id = ? AND user_id != ?", message.RoomID, message.UserID).
			Update("unread_count", gorm.Expr("unread_count + ?", 1)).Error
	})
}

func (s *ChatService) GetMessages(roomID uint, page, limit int) ([]models.ChatMessage, int64, error) {
	var messages []models.ChatMessage
	var total int64

	query := s.db.Model(&models.ChatMessage{}).Where("room_id = ? AND is_deleted = ?", roomID, false)
	
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.Preload("User").Preload("ReplyTo").Preload("ReplyTo.User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&messages).Error

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, total, err
}

func (s *ChatService) GetMessageByID(messageID uint) (*models.ChatMessage, error) {
	var message models.ChatMessage
	err := s.db.Preload("User").Preload("ReplyTo").First(&message, messageID).Error
	return &message, err
}

func (s *ChatService) UpdateMessage(messageID uint, content string) error {
	now := time.Now()
	return s.db.Model(&models.ChatMessage{}).Where("id = ?", messageID).Updates(map[string]interface{}{
		"message":   content,
		"is_edited": true,
		"edited_at": &now,
	}).Error
}

func (s *ChatService) DeleteMessage(messageID uint) error {
	return s.db.Model(&models.ChatMessage{}).Where("id = ?", messageID).Update("is_deleted", true).Error
}

func (s *ChatService) SearchMessages(roomID uint, query string, page, limit int) ([]models.ChatMessage, int64, error) {
	var messages []models.ChatMessage
	var total int64

	dbQuery := s.db.Model(&models.ChatMessage{}).
		Where("room_id = ? AND is_deleted = ? AND message ILIKE ?", roomID, false, "%"+query+"%")
	
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := dbQuery.Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(limit).Find(&messages).Error

	return messages, total, err
}

// Reaction Management
func (s *ChatService) AddReaction(reaction *models.ChatReaction) error {
	// Check if reaction already exists
	var existing models.ChatReaction
	err := s.db.Where("message_id = ? AND user_id = ? AND emoji = ?", 
		reaction.MessageID, reaction.UserID, reaction.Emoji).First(&existing).Error
	
	if err == nil {
		// Reaction already exists
		return errors.New("reaction already exists")
	}
	
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.db.Create(reaction).Error
}

func (s *ChatService) RemoveReaction(messageID, userID uint, emoji string) error {
	return s.db.Where("message_id = ? AND user_id = ? AND emoji = ?", messageID, userID, emoji).
		Delete(&models.ChatReaction{}).Error
}

func (s *ChatService) GetMessageReactions(messageID uint) ([]models.ChatReaction, error) {
	var reactions []models.ChatReaction
	err := s.db.Where("message_id = ?", messageID).Preload("User").Find(&reactions).Error
	return reactions, err
}

// Statistics
func (s *ChatService) GetRoomStats(roomID uint) (map[string]interface{}, error) {
	var messageCount int64
	var memberCount int64

	if err := s.db.Model(&models.ChatMessage{}).Where("room_id = ? AND is_deleted = ?", roomID, false).Count(&messageCount).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.ChatMember{}).Where("room_id = ?", roomID).Count(&memberCount).Error; err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message_count": messageCount,
		"member_count":  memberCount,
	}, nil
}

func (s *ChatService) GetUnreadCount(userID uint) (int64, error) {
	var total int64
	err := s.db.Model(&models.ChatMember{}).
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(unread_count), 0)").
		Scan(&total).Error
	return total, err
}
