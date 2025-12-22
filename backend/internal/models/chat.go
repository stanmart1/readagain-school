package models

import "time"

type ChatRoom struct {
	BaseModel
	Type        string  `gorm:"not null;index" json:"type"` // group, direct, book_discussion
	Name        string  `gorm:"not null" json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	GroupID     *uint   `gorm:"index" json:"group_id"`
	Group       *Group  `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	BookID      *uint   `gorm:"index" json:"book_id"`
	Book        *Book   `gorm:"foreignKey:BookID" json:"book,omitempty"`
	CreatedBy   uint    `gorm:"not null;index" json:"created_by"`
	Creator     *User   `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	IsActive    bool    `gorm:"default:true;index" json:"is_active"`
	LastMessage *string `json:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at"`
}

type ChatMessage struct {
	BaseModel
	RoomID      uint      `gorm:"not null;index" json:"room_id"`
	Room        *ChatRoom `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Message     string    `gorm:"type:text;not null" json:"message"`
	MessageType string    `gorm:"default:text" json:"message_type"` // text, image, file, system
	FileURL     string    `json:"file_url"`
	FileName    string    `json:"file_name"`
	ReplyToID   *uint     `gorm:"index" json:"reply_to_id"`
	ReplyTo     *ChatMessage `gorm:"foreignKey:ReplyToID" json:"reply_to,omitempty"`
	IsEdited    bool      `gorm:"default:false" json:"is_edited"`
	EditedAt    *time.Time `json:"edited_at"`
	IsDeleted   bool      `gorm:"default:false;index" json:"is_deleted"`
}

type ChatMember struct {
	BaseModel
	RoomID       uint       `gorm:"not null;index" json:"room_id"`
	Room         *ChatRoom  `gorm:"foreignKey:RoomID" json:"room,omitempty"`
	UserID       uint       `gorm:"not null;index" json:"user_id"`
	User         *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role         string     `gorm:"default:member" json:"role"` // admin, moderator, member
	LastReadAt   *time.Time `json:"last_read_at"`
	UnreadCount  int        `gorm:"default:0" json:"unread_count"`
	IsMuted      bool       `gorm:"default:false" json:"is_muted"`
	JoinedAt     time.Time  `gorm:"not null" json:"joined_at"`
}

type ChatReaction struct {
	BaseModel
	MessageID uint         `gorm:"not null;index" json:"message_id"`
	Message   *ChatMessage `gorm:"foreignKey:MessageID" json:"message,omitempty"`
	UserID    uint         `gorm:"not null;index" json:"user_id"`
	User      *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Emoji     string       `gorm:"not null" json:"emoji"`
}

type TypingIndicator struct {
	RoomID    uint      `json:"room_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Timestamp time.Time `json:"timestamp"`
}
