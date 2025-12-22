package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{db: db}
}

func (s *CategoryService) ListCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := s.db.Order("name ASC").Find(&categories).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to fetch categories", err)
	}
	return categories, nil
}

func (s *CategoryService) GetCategoryByID(categoryID uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		return nil, utils.NewNotFoundError("Category not found")
	}
	return &category, nil
}

func (s *CategoryService) CreateCategory(name, description string) (*models.Category, error) {
	var existing models.Category
	if err := s.db.Where("name = ?", name).First(&existing).Error; err == nil {
		return nil, utils.NewBadRequestError("Category with this name already exists")
	}

	category := models.Category{
		Name:        name,
		Description: description,
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to create category", err)
	}

	return &category, nil
}

func (s *CategoryService) UpdateCategory(categoryID uint, name, description string) (*models.Category, error) {
	var category models.Category
	if err := s.db.First(&category, categoryID).Error; err != nil {
		return nil, utils.NewNotFoundError("Category not found")
	}

	if name != "" && name != category.Name {
		var existing models.Category
		if err := s.db.Where("name = ? AND id != ?", name, categoryID).First(&existing).Error; err == nil {
			return nil, utils.NewBadRequestError("Category with this name already exists")
		}
		category.Name = name
	}

	if description != "" {
		category.Description = description
	}

	if err := s.db.Save(&category).Error; err != nil {
		return nil, utils.NewInternalServerError("Failed to update category", err)
	}

	return &category, nil
}

func (s *CategoryService) DeleteCategory(categoryID uint) error {
	var booksCount int64
	if err := s.db.Model(&models.Book{}).Where("category_id = ?", categoryID).Count(&booksCount).Error; err != nil {
		return utils.NewInternalServerError("Failed to check category usage", err)
	}

	if booksCount > 0 {
		return utils.NewBadRequestError("Cannot delete category that has books assigned to it")
	}

	if err := s.db.Delete(&models.Category{}, categoryID).Error; err != nil {
		return utils.NewInternalServerError("Failed to delete category", err)
	}

	return nil
}
