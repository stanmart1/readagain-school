package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type ReadingSessionService struct {
	db *gorm.DB
}

func NewReadingSessionService(db *gorm.DB) *ReadingSessionService {
	return &ReadingSessionService{db: db}
}

func (s *ReadingSessionService) GetDB() *gorm.DB {
	return s.db
}

func (s *ReadingSessionService) StartSession(userID, bookID uint) (*models.ReadingSession, error) {
	var library models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&library).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found in library")
	}

	session := models.ReadingSession{
		UserID: userID,
		BookID: bookID,
	}

	if err := s.db.Create(&session).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to start session", err)
	}

	return &session, nil
}

func (s *ReadingSessionService) EndSession(sessionID uint, duration int) error {
	var session models.ReadingSession
	if err := s.db.First(&session, sessionID).Error; err != nil {
		return utils.NewNotFoundError("Session not found")
	}

	session.Duration = duration

	if err := s.db.Save(&session).Error; err != nil {
		return utils.NewInternalServerError("Failed to end session", err)
	}

	return nil
}

func (s *ReadingSessionService) GetUserSessions(userID uint, page, limit int) ([]models.ReadingSession, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.ReadingSession{}).
		Where("user_id = ?", userID).
		Preload("Book")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to count sessions", err)
	}

	var sessions []models.ReadingSession
	if err := query.Scopes(utils.Paginate(params)).Order("created_at DESC").Find(&sessions).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to fetch sessions", err)
	}

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return sessions, &meta, nil
}
