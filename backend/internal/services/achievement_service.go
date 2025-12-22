package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type AchievementService struct {
	db *gorm.DB
}

func NewAchievementService(db *gorm.DB) *AchievementService {
	return &AchievementService{db: db}
}

func (s *AchievementService) GetAllAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement
	if err := s.db.Order("category ASC, target ASC").Find(&achievements).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch achievements", err)
	}
	return achievements, nil
}

func (s *AchievementService) GetUserAchievements(userID uint) ([]models.UserAchievement, error) {
	var userAchievements []models.UserAchievement
	if err := s.db.Where("user_id = ?", userID).
		Preload("Achievement").
		Order("unlocked_at DESC").
		Find(&userAchievements).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch user achievements", err)
	}
	return userAchievements, nil
}

func (s *AchievementService) CheckAndUnlockAchievements(userID uint) ([]models.Achievement, error) {
	var unlockedAchievements []models.Achievement

	var totalBooks int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ?", userID).Count(&totalBooks)

	var completedBooks int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ? AND completed_at IS NOT NULL", userID).Count(&completedBooks)

	var totalReadingTime int64
	s.db.Model(&models.ReadingSession{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration), 0)").Scan(&totalReadingTime)

	var totalSessions int64
	s.db.Model(&models.ReadingSession{}).Where("user_id = ?", userID).Count(&totalSessions)

	criteria := map[string]int64{
		"books_purchased": totalBooks,
		"books_completed": completedBooks,
		"reading_minutes": totalReadingTime,
		"reading_sessions": totalSessions,
	}

	var achievements []models.Achievement
	s.db.Find(&achievements)

	for _, achievement := range achievements {
		if currentValue, exists := criteria[achievement.Type]; exists {
			if int(currentValue) >= achievement.Target {
				var existing models.UserAchievement
				if err := s.db.Where("user_id = ? AND achievement_id = ?", userID, achievement.ID).First(&existing).Error; err != nil {
					userAchievement := models.UserAchievement{
						UserID:        userID,
						AchievementID: achievement.ID,
						Progress:      achievement.Target,
						IsUnlocked:    true,
					}
					if err := s.db.Create(&userAchievement).Error; err == nil {
						unlockedAchievements = append(unlockedAchievements, achievement)
						utils.InfoLogger.Printf("User %d unlocked achievement: %s", userID, achievement.Name)
					}
				}
			}
		}
	}

	return unlockedAchievements, nil
}

func (s *AchievementService) SeedAchievements() error {
	utils.InfoLogger.Println("Achievement seeding skipped - manage via admin panel")
	return nil
}

func (s *AchievementService) CreateAchievement(name, description, icon, achievementType string, target, points int) (*models.Achievement, error) {
	var existing models.Achievement
	if err := s.db.Where("name = ?", name).First(&existing).Error; err == nil {
		return nil, utils.NewBadRequestError("Achievement with this name already exists")
	}

	achievement := models.Achievement{
		Name:        name,
		Description: description,
		Icon:        icon,
		Type:        achievementType,
		Target:      target,
		Points:      points,
	}

	if err := s.db.Create(&achievement).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create achievement", err)
	}

	return &achievement, nil
}

func (s *AchievementService) UpdateAchievement(achievementID uint, updates map[string]interface{}) (*models.Achievement, error) {
	var achievement models.Achievement
	if err := s.db.First(&achievement, achievementID).Error; err != nil {
		return nil, utils.NewNotFoundError("Achievement not found")
	}

	if err := s.db.Model(&achievement).Updates(updates).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update achievement", err)
	}

	return &achievement, nil
}

func (s *AchievementService) DeleteAchievement(achievementID uint) error {
	var userAchievements int64
	s.db.Model(&models.UserAchievement{}).Where("achievement_id = ?", achievementID).Count(&userAchievements)

	if userAchievements > 0 {
		return utils.NewBadRequestError("Cannot delete achievement that users have unlocked")
	}

	if err := s.db.Delete(&models.Achievement{}, achievementID).Error; err != nil {
		return utils.NewInternalServerError("Failed to delete achievement", err)
	}

	return nil
}

