package migrations

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

func RunMigrations(db *gorm.DB) error {
	utils.InfoLogger.Println("Running database migrations...")

	err := db.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.AuthLog{},
		&models.TokenBlacklist{},
		&models.Author{},
		&models.Book{},
		&models.Category{},
		&models.Cart{},
		&models.Order{},
		&models.OrderItem{},
		&models.UserLibrary{},
		&models.ReadingSession{},
		&models.ReadingGoal{},
		&models.Blog{},
		&models.FAQ{},
		&models.Review{},
		&models.SystemSettings{},
		&models.AuditLog{},
		&models.Notification{},
		&models.Achievement{},
		&models.UserAchievement{},
	)

	if err != nil {
		utils.ErrorLogger.Printf("Migration failed: %v", err)
		return err
	}

	utils.InfoLogger.Println("✅ Migrations completed successfully")
	return nil
}
