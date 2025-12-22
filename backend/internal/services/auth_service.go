package services

import (
	"time"

	"gorm.io/gorm"

	"readagain/internal/config"
	"readagain/internal/models"
	"readagain/internal/utils"
)

type AuthService struct {
	db           *gorm.DB
	cfg          *config.Config
	emailService *EmailService
}

func NewAuthService(db *gorm.DB, cfg *config.Config, emailService *EmailService) *AuthService {
	return &AuthService{
		db:           db,
		cfg:          cfg,
		emailService: emailService,
	}
}

func (s *AuthService) Register(email, username, password, firstName, lastName string) (*models.User, error) {
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", email, username).First(&existingUser).Error; err == nil {
		return nil, utils.NewBadRequestError("Email or username already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to hash password", err)
	}

	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPassword,
		FirstName:    firstName,
		LastName:     lastName,
		RoleID:       4,
		IsActive:     false,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create user", err)
	}

	name := firstName
	if name == "" {
		name = username
	}
	go s.emailService.SendWelcomeEmail(email, name)

	return user, nil
}

func (s *AuthService) Login(emailOrUsername, password, ipAddress, userAgent string) (string, string, *models.User, error) {
	var user models.User
	if err := s.db.Preload("Role").Where("email = ? OR username = ?", emailOrUsername, emailOrUsername).First(&user).Error; err != nil {
		s.logAuthAttempt(0, "login", ipAddress, userAgent, false)
		return "", "", nil, utils.NewUnauthorizedError("Invalid credentials")
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		s.logAuthAttempt(user.ID, "login", ipAddress, userAgent, false)
		return "", "", nil, utils.NewUnauthorizedError("Invalid credentials")
	}

	if !user.IsActive {
		return "", "", nil, utils.NewForbiddenError("Account is not active")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.RoleID, s.cfg.JWT.Secret, s.cfg.JWT.ExpireHours)
	if err != nil {
		return "", "", nil, utils.NewInternalServerError("Failed to generate access token", err)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, s.cfg.JWT.Secret, s.cfg.JWT.RefreshExpireDays)
	if err != nil {
		return "", "", nil, utils.NewInternalServerError("Failed to generate refresh token", err)
	}

	now := time.Now()
	user.LastLogin = &now
	s.db.Save(&user)

	s.logAuthAttempt(user.ID, "login", ipAddress, userAgent, true)

	return accessToken, refreshToken, &user, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.ValidateRefreshToken(refreshToken, s.cfg.JWT.Secret)
	if err != nil {
		return "", utils.NewUnauthorizedError("Invalid refresh token")
	}

	var user models.User
	if err := s.db.First(&user, claims.UserID).Error; err != nil {
		return "", utils.NewUnauthorizedError("User not found")
	}

	if !user.IsActive {
		return "", utils.NewForbiddenError("Account is not active")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.RoleID, s.cfg.JWT.Secret, s.cfg.JWT.ExpireHours)
	if err != nil {
		return "", utils.NewInternalServerError("Failed to generate access token", err)
	}

	return accessToken, nil
}

func (s *AuthService) Logout(token string, userID uint) error {
	claims, err := utils.ValidateAccessToken(token, s.cfg.JWT.Secret)
	if err != nil {
		return utils.NewUnauthorizedError("Invalid token")
	}

	blacklist := &models.TokenBlacklist{
		Token:     token,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	if err := s.db.Create(blacklist).Error; err != nil {
		return utils.NewInternalServerError("Failed to blacklist token", err)
	}

	return nil
}

func (s *AuthService) logAuthAttempt(userID uint, action, ipAddress, userAgent string, success bool) {
	if userID == 0 {
		return
	}
	log := &models.AuthLog{
		UserID:    userID,
		Action:    action,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Success:   success,
	}
	s.db.Create(log)
}

func (s *AuthService) RequestPasswordReset(email string) (string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", utils.NewNotFoundError("User not found")
	}

	resetToken := utils.GenerateRandomToken(32)
	expiresAt := time.Now().Add(time.Hour * 1)

	user.VerificationToken = resetToken
	user.VerificationTokenExpires = &expiresAt

	if err := s.db.Save(&user).Error; err != nil {
		return "", utils.NewInternalServerError("Failed to save reset token", err)
	}

	return resetToken, nil
}

func (s *AuthService) ResetPassword(token, newPassword string) error {
	var user models.User
	if err := s.db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return utils.NewNotFoundError("Invalid reset token")
	}

	if user.VerificationTokenExpires == nil || user.VerificationTokenExpires.Before(time.Now()) {
		return utils.NewBadRequestError("Reset token has expired")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.NewInternalServerError("Failed to hash password", err)
	}

	user.PasswordHash = hashedPassword
	user.VerificationToken = ""
	user.VerificationTokenExpires = nil

	if err := s.db.Save(&user).Error; err != nil {
		return utils.NewInternalServerError("Failed to reset password", err)
	}

	return nil
}
