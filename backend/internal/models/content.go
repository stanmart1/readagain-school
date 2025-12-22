package models

import "time"

type Blog struct {
	BaseModel
	Title          string     `gorm:"not null" json:"title" validate:"required"`
	Slug           string     `gorm:"uniqueIndex;not null" json:"slug" validate:"required"`
	Excerpt        string     `gorm:"type:text" json:"excerpt"`
	Content        string     `gorm:"type:text" json:"content"`
	FeaturedImage  string     `json:"featured_image"`
	AuthorID       uint       `json:"author_id"`
	Author         *User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	CategoryID     *uint      `json:"category_id"`
	Status         string     `gorm:"default:draft;index" json:"status"`
	IsFeatured     bool       `gorm:"default:false" json:"is_featured"`
	Views          int        `gorm:"default:0" json:"views"`
	ReadTime       int        `json:"read_time"`
	Tags           string     `json:"tags"`
	PublishedAt    *time.Time `json:"published_at"`
	SEOTitle       string     `json:"seo_title"`
	SEODescription string     `gorm:"type:text" json:"seo_description"`
	SEOKeywords    string     `json:"seo_keywords"`
}

type FAQ struct {
	BaseModel
	Question string `gorm:"not null" json:"question" validate:"required"`
	Answer   string `gorm:"type:text;not null" json:"answer" validate:"required"`
	Category string `json:"category"`
	Order    int    `gorm:"default:0" json:"order"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

type Review struct {
	BaseModel
	UserID     uint   `gorm:"not null;index" json:"user_id"`
	User       *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookID     uint   `gorm:"not null;index" json:"book_id"`
	Book       *Book  `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Rating     int    `gorm:"not null" json:"rating" validate:"required,gte=1,lte=5"`
	Comment    string `gorm:"type:text" json:"comment"`
	Status     string `gorm:"default:pending;index" json:"status"`
	IsFeatured bool   `gorm:"default:false" json:"is_featured"`
}

type Testimonial struct {
	BaseModel
	Name      string `gorm:"not null" json:"name" validate:"required"`
	Role      string `json:"role"`
	Company   string `json:"company"`
	Content   string `gorm:"type:text;not null" json:"content" validate:"required"`
	Avatar    string `json:"avatar"`
	Rating    int    `gorm:"default:5" json:"rating" validate:"gte=1,lte=5"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	Order     int    `gorm:"default:0" json:"order"`
}

type ContactMessage struct {
	BaseModel
	Name    string `gorm:"not null" json:"name" validate:"required"`
	Email   string `gorm:"not null" json:"email" validate:"required,email"`
	Subject string `json:"subject"`
	Message string `gorm:"type:text;not null" json:"message" validate:"required"`
	Status  string `gorm:"default:new;index" json:"status"`
	Reply   string `gorm:"type:text" json:"reply"`
}

