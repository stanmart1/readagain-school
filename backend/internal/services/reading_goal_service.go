package services

import (
	"time"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type ReadingGoalService struct {
	db *gorm.DB
}

func NewReadingGoalService(db *gorm.DB) *ReadingGoalService {
	return &ReadingGoalService{db: db}
}

func (s *ReadingGoalService) CreateGoal(userID uint, goalType string, targetValue int, startDate, endDate time.Time) (*models.ReadingGoal, error) {
	goal := models.ReadingGoal{
		UserID:      userID,
		Type:        goalType,
		Target:      targetValue,
		Current:     0,
		StartDate:   startDate,
		EndDate:     endDate,
		IsCompleted: false,
	}

	if err := s.db.Create(&goal).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create goal", err)
	}

	return &goal, nil
}

func (s *ReadingGoalService) GetUserGoals(userID uint) ([]models.ReadingGoal, error) {
	var goals []models.ReadingGoal
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&goals).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch goals", err)
	}

	for i := range goals {
		s.updateGoalProgress(&goals[i])
	}

	return goals, nil
}

func (s *ReadingGoalService) UpdateGoal(goalID, userID uint, targetValue int) (*models.ReadingGoal, error) {
	var goal models.ReadingGoal
	if err := s.db.Where("id = ? AND user_id = ?", goalID, userID).First(&goal).Error; err != nil {
		return nil, utils.NewNotFoundError("Goal not found")
	}

	goal.Target = targetValue

	if err := s.db.Save(&goal).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update goal", err)
	}

	s.updateGoalProgress(&goal)
	return &goal, nil
}

func (s *ReadingGoalService) DeleteGoal(goalID, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", goalID, userID).Delete(&models.ReadingGoal{})
	if result.Error != nil {
		return utils.NewInternalServerError("Failed to delete goal", result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("Goal not found")
	}
	return nil
}

func (s *ReadingGoalService) updateGoalProgress(goal *models.ReadingGoal) {
	var currentValue int64

	switch goal.Type {
	case "books":
		s.db.Model(&models.UserLibrary{}).
			Where("user_id = ? AND completed_at BETWEEN ? AND ?", goal.UserID, goal.StartDate, goal.EndDate).
			Count(&currentValue)
	case "pages":
		s.db.Model(&models.ReadingSession{}).
			Where("user_id = ? AND created_at BETWEEN ? AND ?", goal.UserID, goal.StartDate, goal.EndDate).
			Select("COALESCE(SUM(pages_read), 0)").
			Scan(&currentValue)
	case "minutes":
		s.db.Model(&models.ReadingSession{}).
			Where("user_id = ? AND created_at BETWEEN ? AND ?", goal.UserID, goal.StartDate, goal.EndDate).
			Select("COALESCE(SUM(duration), 0)").
			Scan(&currentValue)
	}

	goal.Current = int(currentValue)

	if int(currentValue) >= goal.Target && !goal.IsCompleted {
		goal.IsCompleted = true
		s.db.Save(goal)
	}
}
