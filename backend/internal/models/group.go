package models

type Group struct {
	BaseModel
	Name        string `gorm:"not null" json:"name" validate:"required"`
	Description string `gorm:"type:text" json:"description"`
	CreatedBy   uint   `gorm:"not null;index" json:"created_by"`
	Creator     *User  `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	MemberCount int    `gorm:"default:0" json:"member_count"`
}

type GroupMember struct {
	BaseModel
	GroupID uint   `gorm:"not null;index" json:"group_id"`
	Group   *Group `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	UserID  uint   `gorm:"not null;index" json:"user_id"`
	User    *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}
