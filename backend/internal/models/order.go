package models

type Order struct {
	BaseModel
	UserID        uint        `gorm:"not null;index" json:"user_id"`
	User          *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AuthorID      uint        `gorm:"not null;index" json:"author_id"`
	Author        *Author     `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	OrderNumber   string      `gorm:"uniqueIndex;not null" json:"order_number"`
	Subtotal      float64     `gorm:"default:0" json:"subtotal"`
	TaxAmount     float64     `gorm:"default:0" json:"tax_amount"`
	TotalAmount   float64     `gorm:"not null" json:"total_amount"`
	Status        string      `gorm:"default:pending;index" json:"status"`
	PaymentMethod string      `json:"payment_method"`
	Notes         string      `gorm:"type:text" json:"notes"`
	Items         []OrderItem `gorm:"foreignKey:OrderID" json:"items,omitempty"`
}

type OrderItem struct {
	BaseModel
	OrderID   uint    `gorm:"not null;index" json:"order_id"`
	BookID    uint    `gorm:"not null;index" json:"book_id"`
	Book      *Book   `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Quantity  int     `gorm:"not null" json:"quantity" validate:"required,gte=1"`
	Price     float64 `gorm:"not null" json:"price"`
	BookTitle string  `json:"book_title"`
}
