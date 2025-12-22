package services

import (
	"time"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type LibraryService struct {
	db *gorm.DB
}

func NewLibraryService(db *gorm.DB) *LibraryService {
	return &LibraryService{db: db}
}

func (s *LibraryService) GetUserLibrary(userID uint, page, limit int, search string) ([]models.UserLibrary, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.UserLibrary{}).
		Where("user_id = ?", userID).
		Preload("Book.Author.User").
		Preload("Book.Category")

	if search != "" {
		query = query.Joins("JOIN books ON books.id = user_libraries.book_id").
			Where("books.title ILIKE ?", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to count library items", err)
	}

	var library []models.UserLibrary
	if err := query.Scopes(utils.Paginate(params)).Order("user_libraries.created_at DESC").Find(&library).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to fetch library", err)
	}

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return library, &meta, nil
}

func (s *LibraryService) GetLibraryItem(userID, bookID uint) (*models.UserLibrary, error) {
	var library models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).
		Preload("Book.Author.User").
		Preload("Book.Category").
		First(&library).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found in library")
	}
	return &library, nil
}

func (s *LibraryService) UpdateProgress(userID, bookID uint, currentPage, totalPages int, progress float64) error {
	var library models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&library).Error; err != nil {
		return utils.NewNotFoundError("Book not found in library")
	}

	library.CurrentPage = currentPage
	library.Progress = progress

	if progress >= 100 && library.CompletedAt == nil {
		now := time.Now()
		library.CompletedAt = &now
	}

	if err := s.db.Save(&library).Error; err != nil {
		return utils.NewInternalServerError("Failed to update progress", err)
	}

	return nil
}

func (s *LibraryService) GetReadingStatistics(userID uint) (map[string]interface{}, error) {
	var totalBooks int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ?", userID).Count(&totalBooks)

	var completedBooks int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ? AND completed_at IS NOT NULL", userID).Count(&completedBooks)

	var totalReadingTime int64
	s.db.Model(&models.ReadingSession{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration), 0)").Scan(&totalReadingTime)

	var currentStreak int
	var lastReadDate *time.Time
	s.db.Model(&models.ReadingSession{}).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(1).
		Pluck("DATE(created_at)", &lastReadDate)

	if lastReadDate != nil {
		daysSince := int(time.Since(*lastReadDate).Hours() / 24)
		if daysSince <= 1 {
			s.db.Raw(`
				SELECT COUNT(DISTINCT DATE(created_at))
				FROM reading_sessions
				WHERE user_id = ?
				AND created_at >= NOW() - INTERVAL '30 days'
			`, userID).Scan(&currentStreak)
		}
	}

	return map[string]interface{}{
		"total_books":        totalBooks,
		"completed_books":    completedBooks,
		"total_reading_time": totalReadingTime,
		"current_streak":     currentStreak,
	}, nil
}

func (s *LibraryService) GetDashboardStats(userID uint) (map[string]interface{}, error) {
	var booksRead int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ? AND completed_at IS NOT NULL", userID).Count(&booksRead)

	var totalOrders int64
	s.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&totalOrders)

	var reviewsCount int64
	s.db.Model(&models.Review{}).Where("user_id = ?", userID).Count(&reviewsCount)

	return map[string]interface{}{
		"books_read":     booksRead,
		"total_orders":   totalOrders,
		"reviews_count":  reviewsCount,
	}, nil
}

func (s *LibraryService) GetUserAnalytics(userID uint) (map[string]interface{}, error) {
	var totalBooks int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ?", userID).Count(&totalBooks)

	var completedBooks int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ? AND completed_at IS NOT NULL", userID).Count(&completedBooks)

	var totalReadingTime int64
	s.db.Model(&models.ReadingSession{}).Where("user_id = ?", userID).Select("COALESCE(SUM(duration), 0)").Scan(&totalReadingTime)

	var totalPagesRead int64
	s.db.Model(&models.ReadingSession{}).Where("user_id = ?", userID).Select("COALESCE(SUM(pages_read), 0)").Scan(&totalPagesRead)

	return map[string]interface{}{
		"total_books":        totalBooks,
		"completed_books":    completedBooks,
		"total_reading_time": totalReadingTime,
		"total_pages_read":   totalPagesRead,
		"current_streak":     0,
	}, nil
}
