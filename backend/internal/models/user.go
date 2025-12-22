package models

import "time"

type User struct {
	BaseModel
	Email                    string     `gorm:"uniqueIndex;not null" json:"email" validate:"required,email"`
	Username                 string     `gorm:"uniqueIndex;not null" json:"username" validate:"required"`
	PasswordHash             string     `gorm:"not null" json:"-"`
	FirstName                string     `json:"first_name"`
	LastName                 string     `json:"last_name"`
	PhoneNumber              string     `json:"phone_number"`
	SchoolName               string     `json:"school_name"`
	SchoolCategory           string     `json:"school_category"`
	ClassLevel               string     `json:"class_level"`
	Department               string     `json:"department"`
	RoleID                   uint       `gorm:"index" json:"role_id"`
	Role                     *Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	IsActive                 bool       `gorm:"default:false;index" json:"is_active"`
	IsEmailVerified          bool       `gorm:"default:false;index" json:"is_email_verified"`
	VerificationToken        string     `gorm:"index" json:"-"`
	VerificationTokenExpires *time.Time `json:"-"`
	LastLogin                *time.Time `json:"last_login"`
}

type Role struct {
	BaseModel
	Name        string       `gorm:"uniqueIndex;not null" json:"name"` // customer, author, admin, super_admin
	Description string       `json:"description"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

type Permission struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

type AuthLog struct {
	BaseModel
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	User      *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action    string `gorm:"not null" json:"action"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Success   bool   `json:"success"`
}

type TokenBlacklist struct {
	BaseModel
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
}
