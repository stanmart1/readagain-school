package services

import (
	"time"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type BookService struct {
	db *gorm.DB
}

func NewBookService(db *gorm.DB) *BookService {
	return &BookService{db: db}
}

func (s *BookService) GetStats() (map[string]interface{}, error) {
	var totalBooks int64
	var publishedBooks int64
	var featuredBooks int64
	var totalViews int64

	if err := s.db.Model(&models.Book{}).Count(&totalBooks).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Book{}).Where("status = ?", "published").Count(&publishedBooks).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Book{}).Where("is_featured = ?", true).Count(&featuredBooks).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Book{}).Select("COALESCE(SUM(view_count), 0)").Scan(&totalViews).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_books":     totalBooks,
		"published_books": publishedBooks,
		"featured_books":  featuredBooks,
		"total_views":     totalViews,
	}

	return stats, nil
}

type BookFilters struct {
	Search     string
	CategoryID uint
	AuthorID   uint
	MinPrice   float64
	MaxPrice   float64
	IsFeatured *bool
	Status     string
	SortBy     string
	SortOrder  string
}

func (s *BookService) ListBooks(page, limit int, filters BookFilters) ([]models.Book, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.Book{}).
		Preload("Category").
		Preload("Author").
		Preload("Author.User")

	if filters.Search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+filters.Search+"%", "%"+filters.Search+"%")
	}

	if filters.CategoryID > 0 {
		query = query.Where("category_id = ?", filters.CategoryID)
	}

	if filters.AuthorID > 0 {
		query = query.Where("author_id = ?", filters.AuthorID)
	}

	if filters.MinPrice > 0 {
		query = query.Where("price >= ?", filters.MinPrice)
	}

	if filters.MaxPrice > 0 {
		query = query.Where("price <= ?", filters.MaxPrice)
	}

	if filters.IsFeatured != nil {
		query = query.Where("is_featured = ?", *filters.IsFeatured)
	}

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	sortBy := "created_at"
	if filters.SortBy != "" {
		sortBy = filters.SortBy
	}

	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query = query.Order(sortBy + " " + sortOrder)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to count books", err)
	}

	var books []models.Book
	if err := query.Scopes(utils.Paginate(params)).Find(&books).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to fetch books", err)
	}

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return books, &meta, nil
}

func (s *BookService) GetBookByID(bookID uint) (*models.Book, error) {
	var book models.Book
	if err := s.db.Preload("Category").Preload("Author").Preload("Author.User").First(&book, bookID).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found")
	}

	book.ViewCount++
	s.db.Model(&book).Update("view_count", book.ViewCount)

	return &book, nil
}

func (s *BookService) CreateBook(authorID uint, title, description, isbn string, categoryID uint, price float64, coverImage, fileURL string, fileSize int64, pageCount int, status string) (*models.Book, error) {
	var catID *uint
	if categoryID > 0 {
		catID = &categoryID
	}

	book := models.Book{
		AuthorID:    authorID,
		Title:       title,
		Description: description,
		ISBN:        isbn,
		CategoryID:  catID,
		Price:       price,
		CoverImage:  coverImage,
		FilePath:    fileURL,
		FileSize:    fileSize,
		Pages:       pageCount,
		Status:      status,
	}

	if status == "published" {
		now := time.Now()
		book.PublicationDate = &now
	}

	if err := s.db.Create(&book).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create book", err)
	}

	if err := s.db.Preload("Category").Preload("Author").Preload("Author.User").First(&book, book.ID).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found")
	}

	return &book, nil
}

func (s *BookService) UpdateBook(bookID uint, updates map[string]interface{}) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, bookID).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found")
	}

	if status, ok := updates["status"].(string); ok && status == "published" && book.PublicationDate == nil {
		now := time.Now()
		updates["publication_date"] = now
	}

	if err := s.db.Model(&book).Updates(updates).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update book", err)
	}

	if err := s.db.Preload("Category").Preload("Author").Preload("Author.User").First(&book, bookID).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found")
	}

	return &book, nil
}

func (s *BookService) DeleteBook(bookID uint) error {
	if err := s.db.Delete(&models.Book{}, bookID).Error; err != nil {
		return utils.NewInternalServerError("Failed to delete book", err)
	}
	return nil
}

func (s *BookService) GetFeaturedBooks(limit int) ([]models.Book, error) {
	var books []models.Book
	if err := s.db.Where("is_featured = ? AND status = ?", true, "published").
		Preload("Category").
		Preload("Author").
		Preload("Author.User").
		Order("created_at DESC").
		Limit(limit).
		Find(&books).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch featured books", err)
	}
	return books, nil
}

func (s *BookService) GetNewReleases(limit int) ([]models.Book, error) {
	var books []models.Book
	if err := s.db.Where("status = ?", "published").
		Preload("Category").
		Preload("Author").
		Preload("Author.User").
		Order("publication_date DESC").
		Limit(limit).
		Find(&books).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch new releases", err)
	}
	return books, nil
}

func (s *BookService) GetBestsellers(limit int) ([]models.Book, error) {
	var books []models.Book
	if err := s.db.Where("status = ?", "published").
		Preload("Category").
		Preload("Author").
		Preload("Author.User").
		Order("download_count DESC").
		Limit(limit).
		Find(&books).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch bestsellers", err)
	}
	return books, nil
}

func (s *BookService) ToggleFeatured(bookID uint, isFeatured bool) (*models.Book, error) {
	var book models.Book
	if err := s.db.First(&book, bookID).Error; err != nil {
		return nil, utils.NewNotFoundError("Book not found")
	}

	book.IsFeatured = isFeatured
	if err := s.db.Save(&book).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update book", err)
	}

	return &book, nil
}
