package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type AuthorService struct {
	db *gorm.DB
}

func NewAuthorService(db *gorm.DB) *AuthorService {
	return &AuthorService{db: db}
}

func (s *AuthorService) GetStats() (map[string]interface{}, error) {
	var totalAuthors int64
	var activeAuthors int64

	if err := s.db.Model(&models.Author{}).Count(&totalAuthors).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Author{}).Where("status = ?", "active").Count(&activeAuthors).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_authors":  totalAuthors,
		"active_authors": activeAuthors,
	}

	return stats, nil
}

func (s *AuthorService) ListAuthors(page, limit int, search string) ([]map[string]interface{}, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.Author{}).Preload("User")

	if search != "" {
		query = query.Joins("JOIN users ON users.id = authors.user_id").
			Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR authors.business_name ILIKE ?",
				"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to count authors", err)
	}

	var authors []models.Author
	if err := query.Scopes(utils.Paginate(params)).Find(&authors).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to fetch authors", err)
	}

	// Add books_count to each author
	result := make([]map[string]interface{}, len(authors))
	for i, author := range authors {
		var booksCount int64
		s.db.Model(&models.Book{}).Where("author_id = ?", author.ID).Count(&booksCount)
		
		result[i] = map[string]interface{}{
			"id":            author.ID,
			"created_at":    author.CreatedAt,
			"updated_at":    author.UpdatedAt,
			"user_id":       author.UserID,
			"user":          author.User,
			"business_name": author.BusinessName,
			"bio":           author.Bio,
			"website":       author.Website,
			"status":        author.Status,
			"books_count":   booksCount,
		}
	}

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return result, &meta, nil
}

func (s *AuthorService) GetAuthorByID(authorID uint) (*models.Author, error) {
	var author models.Author
	if err := s.db.Preload("User").First(&author, authorID).Error; err != nil {
		return nil, utils.NewNotFoundError("Author not found")
	}
	return &author, nil
}

func (s *AuthorService) GetAuthorByUserID(userID uint) (*models.Author, error) {
	return nil, utils.NewNotFoundError("Author profile not found")
}

func (s *AuthorService) CreateAuthor(penName, bio, website string) (*models.Author, error) {
	author := models.Author{
		BusinessName: penName,
		Bio:          bio,
		Website:      website,
		Status:       "active",
	}

	if err := s.db.Create(&author).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create author", err)
	}

	return &author, nil
}

func (s *AuthorService) UpdateAuthor(authorID uint, updates map[string]interface{}) (*models.Author, error) {
	var author models.Author
	if err := s.db.First(&author, authorID).Error; err != nil {
		return nil, utils.NewNotFoundError("Author not found")
	}

	if err := s.db.Model(&author).Updates(updates).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update author", err)
	}

	if err := s.db.Preload("User").First(&author, authorID).Error; err != nil {
		return nil, utils.NewNotFoundError("Author not found")
	}

	return &author, nil
}

func (s *AuthorService) DeleteAuthor(authorID uint) error {
	var booksCount int64
	if err := s.db.Model(&models.Book{}).Where("author_id = ?", authorID).Count(&booksCount).Error; err != nil {
		return utils.NewInternalServerError("Failed to check author books", err)
	}

	if booksCount > 0 {
		return utils.NewBadRequestError("Cannot delete author who has published books")
	}

	if err := s.db.Delete(&models.Author{}, authorID).Error; err != nil {
		return utils.NewInternalServerError("Failed to delete author", err)
	}

	return nil
}
