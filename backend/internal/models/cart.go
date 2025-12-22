package models

type Cart struct {
	BaseModel
	UserID   uint  `gorm:"not null;index" json:"user_id"`
	User     *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookID   uint  `gorm:"not null;index" json:"book_id"`
	Book     *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Quantity int   `gorm:"not null;default:1" json:"quantity" validate:"required,gte=1"`
}
