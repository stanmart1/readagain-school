package models

import "time"

type Author struct {
	BaseModel
	UserID       uint   `gorm:"uniqueIndex;not null" json:"user_id"`
	User         *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BusinessName string `gorm:"not null" json:"business_name" validate:"required"`
	Bio          string `gorm:"type:text" json:"bio"`
	Photo        string `json:"photo"`
	Website      string `json:"website"`
	Email        string `json:"email" validate:"omitempty,email"`
	Status       string `gorm:"default:active" json:"status"`
}

type Book struct {
	BaseModel
	AuthorID         uint       `gorm:"not null;index" json:"author_id" validate:"required"`
	Author           *Author    `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	Title            string     `gorm:"not null" json:"title" validate:"required"`
	Subtitle         string     `json:"subtitle"`
	Description      string     `gorm:"type:text" json:"description"`
	ShortDescription string     `gorm:"type:text" json:"short_description"`
	Price            float64    `gorm:"not null" json:"price" validate:"required,gte=0"`
	CoverImage       string     `json:"cover_image"`
	FilePath         string     `json:"file_path"`
	FileSize         int64      `json:"file_size"`
	CategoryID       *uint      `json:"category_id"`
	Category         *Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	ISBN             string     `gorm:"uniqueIndex" json:"isbn"`
	IsFeatured       bool       `gorm:"default:false;index" json:"is_featured"`
	IsBestseller     bool       `gorm:"default:false;index" json:"is_bestseller"`
	IsNewRelease     bool       `gorm:"default:false;index" json:"is_new_release"`
	Status           string     `gorm:"default:published;index" json:"status"`
	IsActive         bool       `gorm:"default:true;index" json:"is_active"`
	Pages            int        `json:"pages"`
	Language         string     `gorm:"default:English" json:"language"`
	Publisher        string     `json:"publisher"`
	PublicationDate  *time.Time `json:"publication_date"`
	DownloadCount    int        `gorm:"default:0" json:"download_count"`
	ViewCount        int        `gorm:"default:0" json:"view_count"`
	SEOTitle         string     `json:"seo_title"`
	SEODescription   string     `gorm:"type:text" json:"seo_description"`
	SEOKeywords      string     `json:"seo_keywords"`
}

type Category struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null" json:"name" validate:"required"`
	Description string `gorm:"type:text" json:"description"`
	Status      string `gorm:"default:active" json:"status"`
}
