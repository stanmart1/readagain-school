package models

import "time"

type UserLibrary struct {
	BaseModel
	UserID      uint       `gorm:"not null;index" json:"user_id"`
	User        *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookID      uint       `gorm:"not null;index" json:"book_id"`
	Book        *Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Progress    float64    `gorm:"default:0" json:"progress"`
	CurrentPage int        `gorm:"default:0" json:"current_page"`
	LastReadAt  *time.Time `json:"last_read_at"`
	CompletedAt *time.Time `json:"completed_at"`
	IsFavorite  bool       `gorm:"default:false" json:"is_favorite"`
	Rating      int        `json:"rating" validate:"omitempty,gte=1,lte=5"`
}

type ReadingSession struct {
	BaseModel
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	User      *User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookID    uint       `gorm:"not null;index" json:"book_id"`
	Book      *Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
	StartTime time.Time  `gorm:"not null" json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Duration  int        `json:"duration"`
	PagesRead int        `json:"pages_read"`
	StartPage int        `json:"start_page"`
	EndPage   int        `json:"end_page"`
}

type ReadingGoal struct {
	BaseModel
	UserID      uint      `gorm:"not null;index" json:"user_id"`
	User        *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        string    `gorm:"not null" json:"type" validate:"required"`
	Target      int       `gorm:"not null" json:"target" validate:"required,gte=1"`
	Current     int       `gorm:"default:0" json:"current"`
	StartDate   time.Time `gorm:"not null" json:"start_date"`
	EndDate     time.Time `gorm:"not null" json:"end_date"`
	IsCompleted bool      `gorm:"default:false" json:"is_completed"`
}

type Bookmark struct {
	BaseModel
	UserID   uint   `gorm:"not null;index" json:"user_id"`
	User     *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookID   uint   `gorm:"not null;index" json:"book_id"`
	Book     *Book  `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Page     int    `gorm:"not null" json:"page"`
	Location string `json:"location"`
	Note     string `gorm:"type:text" json:"note"`
}

type Note struct {
	BaseModel
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookID    uint   `gorm:"not null;index" json:"book_id"`
	Book      *Book  `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Page      int    `gorm:"not null" json:"page"`
	Content   string `gorm:"type:text;not null" json:"content"`
	Highlight string `gorm:"type:text" json:"highlight"`
}

