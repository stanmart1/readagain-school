package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type EReaderService struct {
	db *gorm.DB
}

func NewEReaderService(db *gorm.DB) *EReaderService {
	return &EReaderService{db: db}
}

func (s *EReaderService) ValidateBookAccess(userID, bookID uint) (*models.Book, error) {
	var library models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&library).Error; err != nil {
		return nil, utils.NewForbiddenError("You don't own this book")
	}

	var book models.Book
	if err := s.db.Preload("Author").First(&book, bookID).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found")
	}

	return &book, nil
}

func (s *EReaderService) CreateBookmark(userID, bookID uint, page int, location, note string) (*models.Bookmark, error) {
	var library models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&library).Error; err != nil {
		return nil, utils.NewForbiddenError("You don't own this book")
	}

	bookmark := models.Bookmark{
		UserID:   userID,
		BookID:   bookID,
		Page:     page,
		Location: location,
		Note:     note,
	}

	if err := s.db.Create(&bookmark).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create bookmark", err)
	}

	return &bookmark, nil
}

func (s *EReaderService) GetBookmarks(userID, bookID uint) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).
		Order("page ASC").
		Find(&bookmarks).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch bookmarks", err)
	}
	return bookmarks, nil
}

func (s *EReaderService) DeleteBookmark(bookmarkID, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", bookmarkID, userID).Delete(&models.Bookmark{})
	if result.Error != nil {
		return utils.NewInternalServerError("Failed to delete bookmark", result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("Bookmark not found")
	}
	return nil
}

func (s *EReaderService) CreateNote(userID, bookID uint, page int, content, highlight string) (*models.Note, error) {
	var library models.UserLibrary
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).First(&library).Error; err != nil {
		return nil, utils.NewForbiddenError("You don't own this book")
	}

	note := models.Note{
		UserID:    userID,
		BookID:    bookID,
		Page:      page,
		Content:   content,
		Highlight: highlight,
	}

	if err := s.db.Create(&note).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create note", err)
	}

	return &note, nil
}

func (s *EReaderService) GetNotes(userID, bookID uint) ([]models.Note, error) {
	var notes []models.Note
	if err := s.db.Where("user_id = ? AND book_id = ?", userID, bookID).
		Order("page ASC").
		Find(&notes).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch notes", err)
	}
	return notes, nil
}

func (s *EReaderService) UpdateNote(noteID, userID uint, content string) (*models.Note, error) {
	var note models.Note
	if err := s.db.Where("id = ? AND user_id = ?", noteID, userID).First(&note).Error; err != nil {
		return nil, utils.NewNotFoundError("Note not found")
	}

	note.Content = content

	if err := s.db.Save(&note).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update note", err)
	}

	return &note, nil
}

func (s *EReaderService) DeleteNote(noteID, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", noteID, userID).Delete(&models.Note{})
	if result.Error != nil {
		return utils.NewInternalServerError("Failed to delete note", result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.NewNotFoundError("Note not found")
	}
	return nil
}
