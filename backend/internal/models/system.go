package models

type SystemSettings struct {
	BaseModel
	Key         string `gorm:"uniqueIndex;not null" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	DataType    string `json:"data_type"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
}

type AuditLog struct {
	BaseModel
	UserID     uint   `gorm:"index" json:"user_id"`
	User       *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Action     string `gorm:"not null" json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   uint   `json:"entity_id"`
	OldValue   string `gorm:"type:text" json:"old_value"`
	NewValue   string `gorm:"type:text" json:"new_value"`
	IPAddress  string `json:"ip_address"`
	UserAgent  string `json:"user_agent"`
}

type Notification struct {
	BaseModel
	UserID   uint   `gorm:"not null;index" json:"user_id"`
	User     *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type     string `gorm:"not null" json:"type"`
	Title    string `gorm:"not null" json:"title"`
	Message  string `gorm:"type:text" json:"message"`
	IsRead   bool   `gorm:"default:false;index" json:"is_read"`
	ActionURL string `json:"action_url"`
}

type Achievement struct {
	BaseModel
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Icon        string `json:"icon"`
	Type        string `json:"type"`
	Target      int    `json:"target"`
	Points      int    `json:"points"`
}

type UserAchievement struct {
	BaseModel
	UserID        uint         `gorm:"not null;index" json:"user_id"`
	User          *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AchievementID uint         `gorm:"not null;index" json:"achievement_id"`
	Achievement   *Achievement `gorm:"foreignKey:AchievementID" json:"achievement,omitempty"`
	Progress      int          `gorm:"default:0" json:"progress"`
	IsUnlocked    bool         `gorm:"default:false" json:"is_unlocked"`
}

type AboutPage struct {
	BaseModel
	Title       string `json:"title"`
	Content     string `gorm:"type:text" json:"content"`
	Mission     string `gorm:"type:text" json:"mission"`
	Vision      string `gorm:"type:text" json:"vision"`
	TeamSection string `gorm:"type:text" json:"team_section"`
	Values      string `gorm:"type:text" json:"values"`
}

type Activity struct {
	BaseModel
	UserID      uint   `gorm:"not null;index" json:"user_id"`
	User        *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Type        string `gorm:"not null;index" json:"type"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EntityType  string `json:"entity_type"`
	EntityID    uint   `json:"entity_id"`
	Metadata    string `gorm:"type:text" json:"metadata"`
}

type Wishlist struct {
	BaseModel
	UserID uint  `gorm:"not null;index" json:"user_id"`
	User   *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	BookID uint  `gorm:"not null;index" json:"book_id"`
	Book   *Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}



