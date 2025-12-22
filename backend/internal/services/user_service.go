package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}
	return &user, nil
}

func (s *UserService) UpdateProfile(userID uint, firstName, lastName, phoneNumber, schoolName, schoolCategory, classLevel, department string) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.PhoneNumber = phoneNumber
	user.SchoolName = schoolName
	user.SchoolCategory = schoolCategory
	user.ClassLevel = classLevel
	user.Department = department

	if err := s.db.Save(&user).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update profile", err)
	}

	return &user, nil
}

func (s *UserService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return utils.NewNotFoundError("User not found")
	}

	if !utils.CheckPassword(currentPassword, user.PasswordHash) {
		return utils.NewUnauthorizedError("Current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.NewInternalServerError("Failed to hash password", err)
	}

	user.PasswordHash = hashedPassword
	if err := s.db.Save(&user).Error; err != nil {
		return utils.NewInternalServerError("Failed to change password", err)
	}

	return nil
}

func (s *UserService) DeactivateUser(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return utils.NewNotFoundError("User not found")
	}

	user.IsActive = false
	if err := s.db.Save(&user).Error; err != nil {
		return utils.NewInternalServerError("Failed to deactivate user", err)
	}

	return nil
}

func (s *UserService) ReactivateUser(userID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return utils.NewNotFoundError("User not found")
	}

	user.IsActive = true
	if err := s.db.Save(&user).Error; err != nil {
		return utils.NewInternalServerError("Failed to reactivate user", err)
	}

	return nil
}

func (s *UserService) DeleteUser(userID uint) error {
	if err := s.db.Unscoped().Delete(&models.User{}, userID).Error; err != nil {
		return utils.NewInternalServerError("Failed to delete user", err)
	}
	return nil
}

func (s *UserService) ListUsers(page, limit int, search string) ([]models.User, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.User{}).Preload("Role")

	if search != "" {
		query = query.Where("email ILIKE ? OR username ILIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to count users", err)
	}

	var users []models.User
	if err := query.Scopes(utils.Paginate(params)).Find(&users).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to fetch users", err)
	}

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return users, &meta, nil
}

func (s *UserService) UpdateUser(userID uint, updates map[string]interface{}) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}

	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update user", err)
	}

	if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}

	return &user, nil
}

func (s *UserService) AssignRole(userID, roleID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return utils.NewNotFoundError("User not found")
	}

	var role models.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return utils.NewNotFoundError("Role not found")
	}

	user.RoleID = roleID
	if err := s.db.Save(&user).Error; err != nil {
		return utils.NewInternalServerError("Failed to assign role", err)
	}

	return nil
}

func (s *UserService) RemoveRole(userID, roleID uint) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return utils.NewNotFoundError("User not found")
	}

	if user.RoleID != roleID {
		return utils.NewBadRequestError("User does not have this role")
	}

	user.RoleID = 0
	if err := s.db.Save(&user).Error; err != nil {
		return utils.NewInternalServerError("Failed to remove role", err)
	}

	return nil
}

func (s *UserService) GetUserStats(userID uint) (map[string]interface{}, error) {
	var booksRead int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ? AND completed_at IS NOT NULL", userID).Count(&booksRead)

	var totalReadingTime int64
	s.db.Model(&models.ReadingSession{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration), 0)").Scan(&totalReadingTime)

	var activeGoals int64
	s.db.Model(&models.ReadingGoal{}).Where("user_id = ? AND is_completed = ?", userID, false).Count(&activeGoals)

	return map[string]interface{}{
		"books_read":         booksRead,
		"total_reading_time": totalReadingTime,
		"active_goals":       activeGoals,
	}, nil
}

func (s *UserService) CreateUserByAdmin(email, username, password, firstName, lastName string, roleID uint, isActive bool) (*models.User, error) {
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", email, username).First(&existingUser).Error; err == nil {
		return nil, utils.NewBadRequestError("User with this email or username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to hash password", err)
	}

	user := models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
		RoleID:       roleID,
		IsActive:     isActive,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create user", err)
	}

	if err := s.db.Preload("Role").First(&user, user.ID).Error; err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}

	return &user, nil
}

func (s *UserService) ToggleUserStatus(userID uint, isActive bool) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}

	user.IsActive = isActive
	if err := s.db.Save(&user).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update user status", err)
	}

	if err := s.db.Preload("Role").First(&user, userID).Error; err != nil {
		return nil, utils.NewNotFoundError("User not found")
	}

	return &user, nil
}

func (s *UserService) AdminResetPassword(userID uint, newPassword string) error {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return utils.NewNotFoundError("User not found")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.NewInternalServerError("Failed to hash password", err)
	}

	user.PasswordHash = hashedPassword
	if err := s.db.Save(&user).Error; err != nil {
		return utils.NewInternalServerError("Failed to reset password", err)
	}

	return nil
}

func (s *UserService) BulkActivate(userIDs []uint) error {
	if len(userIDs) == 0 {
		return utils.NewBadRequestError("No user IDs provided")
	}

	if err := s.db.Model(&models.User{}).Where("id IN ?", userIDs).Update("is_active", true).Error; err != nil {
		return utils.NewInternalServerError("Failed to activate users", err)
	}

	return nil
}

func (s *UserService) BulkDeactivate(userIDs []uint) error {
	if len(userIDs) == 0 {
		return utils.NewBadRequestError("No user IDs provided")
	}

	if err := s.db.Model(&models.User{}).Where("id IN ?", userIDs).Update("is_active", false).Error; err != nil {
		return utils.NewInternalServerError("Failed to deactivate users", err)
	}

	return nil
}

func (s *UserService) BulkDelete(userIDs []uint) error {
	if len(userIDs) == 0 {
		return utils.NewBadRequestError("No user IDs provided")
	}

	if err := s.db.Unscoped().Delete(&models.User{}, userIDs).Error; err != nil {
		return utils.NewInternalServerError("Failed to delete users", err)
	}

	return nil
}


