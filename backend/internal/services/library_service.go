package services

import (
	"fmt"
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

	var booksInLibrary int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ?", userID).Count(&booksInLibrary)

	var reviewsCount int64
	s.db.Model(&models.Review{}).Where("user_id = ?", userID).Count(&reviewsCount)

	return map[string]interface{}{
		"books_read":       booksRead,
		"books_in_library": booksInLibrary,
		"reviews_count":    reviewsCount,
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


// Admin methods
func (s *LibraryService) GetAllAssignments(skip, limit int, search, status, userID, sortBy, sortOrder string) ([]map[string]interface{}, int, error) {
	query := s.db.Table("user_libraries ul").
		Select(`ul.id, ul.user_id, ul.book_id, ul.progress, ul.created_at as assigned_at,
			CONCAT(u.first_name, ' ', u.last_name) as user_name, u.email as user_email,
			b.title as book_title, a.business_name as book_author,
			CASE 
				WHEN ul.progress = 0 THEN 'unread'
				WHEN ul.progress = 100 THEN 'completed'
				ELSE 'reading'
			END as status`).
		Joins("JOIN users u ON ul.user_id = u.id").
		Joins("JOIN books b ON ul.book_id = b.id").
		Joins("LEFT JOIN authors a ON b.author_id = a.id")

	if search != "" {
		query = query.Where("u.first_name ILIKE ? OR u.last_name ILIKE ? OR u.email ILIKE ? OR b.title ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		if status == "unread" {
			query = query.Where("ul.progress = 0")
		} else if status == "completed" {
			query = query.Where("ul.progress = 100")
		} else if status == "reading" {
			query = query.Where("ul.progress > 0 AND ul.progress < 100")
		}
	}
	if userID != "" {
		query = query.Where("ul.user_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	// Prefix sortBy with table alias if it's created_at
	if sortBy == "created_at" {
		sortBy = "ul.created_at"
	}
	orderClause := sortBy + " " + sortOrder
	query = query.Order(orderClause).Offset(skip).Limit(limit)

	var assignments []map[string]interface{}
	if err := query.Scan(&assignments).Error; err != nil {
		return nil, 0, err
	}

	return assignments, int(total), nil
}

func (s *LibraryService) GetLibraryStats() (map[string]interface{}, error) {
	var totalAssignments, activeReaders int64
	var avgProgress, completionRate float64

	s.db.Model(&models.UserLibrary{}).Count(&totalAssignments)
	s.db.Model(&models.UserLibrary{}).Distinct("user_id").Count(&activeReaders)
	s.db.Model(&models.UserLibrary{}).Select("COALESCE(AVG(progress), 0)").Scan(&avgProgress)

	var completed int64
	s.db.Model(&models.UserLibrary{}).Where("progress = 100").Count(&completed)
	if totalAssignments > 0 {
		completionRate = float64(completed) / float64(totalAssignments) * 100
	}

	return map[string]interface{}{
		"total_assignments": totalAssignments,
		"active_readers":    activeReaders,
		"avg_progress":      avgProgress,
		"completion_rate":   completionRate,
	}, nil
}

func (s *LibraryService) AssignBook(userID, bookID uint, format string) error {
	var exists int64
	s.db.Model(&models.UserLibrary{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&exists)
	if exists > 0 {
		return fmt.Errorf("book already assigned to this user")
	}

	library := &models.UserLibrary{
		UserID: userID,
		BookID: bookID,
	}

	return s.db.Create(library).Error
}

func (s *LibraryService) BulkAssignBook(userIDs []uint, bookID uint, format string) (int, error) {
	count := 0
	for _, userID := range userIDs {
		var exists int64
		s.db.Model(&models.UserLibrary{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&exists)
		if exists == 0 {
			library := &models.UserLibrary{
				UserID: userID,
				BookID: bookID,
			}
			if err := s.db.Create(library).Error; err == nil {
				count++
			}
		}
	}
	return count, nil
}

func (s *LibraryService) BulkRemoveAssignments(assignmentIDs []uint) (int, error) {
	result := s.db.Where("id IN ?", assignmentIDs).Delete(&models.UserLibrary{})
	return int(result.RowsAffected), result.Error
}

func (s *LibraryService) RemoveAssignment(id uint) error {
	return s.db.Delete(&models.UserLibrary{}, id).Error
}

func (s *LibraryService) GetBooksWithStudents(search string) ([]map[string]interface{}, error) {
	query := s.db.Table("books b").
		Select(`b.id, b.title, b.author, b.cover_image,
			COUNT(DISTINCT ul.user_id) as student_count,
			COALESCE(AVG(ul.progress), 0) as avg_progress,
			COUNT(CASE WHEN ul.progress = 100 THEN 1 END) as completed_count`).
		Joins("LEFT JOIN user_libraries ul ON b.id = ul.book_id").
		Group("b.id, b.title, b.author, b.cover_image")

	if search != "" {
		query = query.Where("b.title ILIKE ? OR b.author ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var books []map[string]interface{}
	if err := query.Scan(&books).Error; err != nil {
		return nil, err
	}

	// Get students for each book
	for i, book := range books {
		bookID := book["id"]
		var students []map[string]interface{}
		s.db.Table("user_libraries ul").
			Select(`ul.id as assignment_id, ul.progress, 
				CASE 
					WHEN ul.progress = 0 THEN 'unread'
					WHEN ul.progress = 100 THEN 'completed'
					ELSE 'reading'
				END as status,
				u.id as user_id, u.name, u.email, u.class_level, u.school_name`).
			Joins("JOIN users u ON ul.user_id = u.id").
			Where("ul.book_id = ?", bookID).
			Scan(&students)
		books[i]["students"] = students
	}

	return books, nil
}

func (s *LibraryService) GetAssignmentDetails(id uint) (map[string]interface{}, error) {
	var assignment map[string]interface{}
	err := s.db.Table("user_libraries ul").
		Select(`ul.*, CONCAT(u.first_name, ' ', u.last_name) as user_name, u.email as user_email, 
			b.title as book_title, a.business_name as book_author,
			CASE 
				WHEN ul.progress = 0 THEN 'unread'
				WHEN ul.progress = 100 THEN 'completed'
				ELSE 'reading'
			END as status`).
		Joins("JOIN users u ON ul.user_id = u.id").
		Joins("JOIN books b ON ul.book_id = b.id").
		Joins("LEFT JOIN authors a ON b.author_id = a.id").
		Where("ul.id = ?", id).
		Scan(&assignment).Error

	if err != nil {
		return nil, err
	}

	// Get notes
	var notes []models.Note
	s.db.Where("user_library_id = ?", id).Find(&notes)

	return map[string]interface{}{
		"assignment": assignment,
		"highlights": []interface{}{},
		"notes":      notes,
	}, nil
}

func (s *LibraryService) GetAssignmentAnalytics(id uint) (map[string]interface{}, error) {
	var assignment models.UserLibrary
	if err := s.db.First(&assignment, id).Error; err != nil {
		return nil, err
	}

	var totalSessions, totalReadingTime int64
	s.db.Model(&models.ReadingSession{}).Where("user_id = ? AND book_id = ?", assignment.UserID, assignment.BookID).Count(&totalSessions)
	s.db.Model(&models.ReadingSession{}).Where("user_id = ? AND book_id = ?", assignment.UserID, assignment.BookID).Select("COALESCE(SUM(duration), 0)").Scan(&totalReadingTime)

	avgSessionTime := int64(0)
	if totalSessions > 0 {
		avgSessionTime = totalReadingTime / totalSessions
	}

	// Get goals
	var goals []models.ReadingGoal
	s.db.Where("user_id = ? AND book_id = ?", assignment.UserID, assignment.BookID).Find(&goals)

	return map[string]interface{}{
		"total_sessions":      totalSessions,
		"total_reading_time":  totalReadingTime / 60,
		"avg_session_time":    avgSessionTime / 60,
		"reading_streak":      0,
		"goals":               goals,
	}, nil
}
