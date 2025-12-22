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

func (s *AuthorService) ListAuthors(page, limit int, search string) ([]models.Author, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.Author{}).Preload("User")

	if search != "" {
		query = query.Joins("JOIN users ON users.id = authors.user_id").
			Where("users.first_name ILIKE ? OR users.last_name ILIKE ? OR authors.pen_name ILIKE ?",
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

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return authors, &meta, nil
}

func (s *AuthorService) GetAuthorByID(authorID uint) (*models.Author, error) {
	var author models.Author
	if err := s.db.Preload("User").First(&author, authorID).Error; err != nil {
		return nil, utils.NewNotFoundError("Author not found")
	}
	return &author, nil
}

func (s *AuthorService) GetAuthorByUserID(userID uint) (*models.Author, error) {
	var author models.Author
	if err := s.db.Preload("User").Where("user_id = ?", userID).First(&author).Error; err != nil {
		return nil, utils.NewNotFoundError("Author profile not found")
	}
	return &author, nil
}

func (s *AuthorService) CreateAuthor(userID uint, penName, bio, website string) (*models.Author, error) {
	var existing models.Author
	if err := s.db.Where("user_id = ?", userID).First(&existing).Error; err == nil {
		return nil, utils.NewBadRequestError("Author profile already exists for this user")
	}

	author := models.Author{
		UserID:       userID,
		BusinessName: penName,
		Bio:          bio,
		Website:      website,
	}

	if err := s.db.Create(&author).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create author", err)
	}

	if err := s.db.Preload("User").First(&author, author.ID).Error; err != nil {
		return nil, utils.NewNotFoundError("Author not found")
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
