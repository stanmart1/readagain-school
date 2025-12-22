package models

type Work struct {
	BaseModel
	Title       string `gorm:"not null" json:"title" validate:"required"`
	Description string `gorm:"type:text" json:"description"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	Order       int    `gorm:"default:0" json:"order"`
}
