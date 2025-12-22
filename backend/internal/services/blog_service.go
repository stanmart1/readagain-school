package services

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type BlogService struct {
	db *gorm.DB
}

func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{db: db}
}

func (s *BlogService) GetStats() (map[string]interface{}, error) {
	var totalBlogs int64
	var publishedBlogs int64
	var draftBlogs int64
	var totalViews int64

	if err := s.db.Model(&models.Blog{}).Count(&totalBlogs).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Blog{}).Where("status = ?", "published").Count(&publishedBlogs).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Blog{}).Where("status = ?", "draft").Count(&draftBlogs).Error; err != nil {
		return nil, err
	}

	if err := s.db.Model(&models.Blog{}).Select("COALESCE(SUM(views), 0)").Scan(&totalViews).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_blogs":     totalBlogs,
		"published_blogs": publishedBlogs,
		"draft_blogs":     draftBlogs,
		"total_views":     totalViews,
	}

	return stats, nil
}

func (s *BlogService) List(page, limit int, status, search, tag string) ([]models.Blog, *utils.PaginationMeta, error) {
	var blogs []models.Blog
	var total int64

	query := s.db.Model(&models.Blog{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if search != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if tag != "" {
		query = query.Where("tags ILIKE ?", "%"+tag+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, nil, err
	}

	offset := (page - 1) * limit
	if err := query.Preload("Author").Order("created_at DESC").Offset(offset).Limit(limit).Find(&blogs).Error; err != nil {
		return nil, nil, err
	}

	meta := &utils.PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}

	return blogs, meta, nil
}

func (s *BlogService) GetByID(id uint) (*models.Blog, error) {
	var blog models.Blog
	if err := s.db.Preload("Author").First(&blog, id).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (s *BlogService) GetBySlug(slug string) (*models.Blog, error) {
	var blog models.Blog
	if err := s.db.Preload("Author").Where("slug = ?", slug).First(&blog).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (s *BlogService) Create(blog *models.Blog) error {
	if blog.Slug == "" {
		blog.Slug = s.generateSlug(blog.Title)
	}

	blog.ReadTime = s.calculateReadTime(blog.Content)

	if blog.Status == "published" && blog.PublishedAt == nil {
		now := time.Now()
		blog.PublishedAt = &now
	}

	return s.db.Create(blog).Error
}

func (s *BlogService) Update(id uint, updates map[string]interface{}) error {
	blog, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if title, ok := updates["title"].(string); ok && title != blog.Title {
		updates["slug"] = s.generateSlug(title)
	}

	if content, ok := updates["content"].(string); ok {
		updates["read_time"] = s.calculateReadTime(content)
	}

	if status, ok := updates["status"].(string); ok && status == "published" && blog.PublishedAt == nil {
		updates["published_at"] = time.Now()
	}

	return s.db.Model(&models.Blog{}).Where("id = ?", id).Updates(updates).Error
}

func (s *BlogService) Delete(id uint) error {
	return s.db.Delete(&models.Blog{}, id).Error
}

func (s *BlogService) IncrementViews(id uint) error {
	return s.db.Model(&models.Blog{}).Where("id = ?", id).UpdateColumn("views", gorm.Expr("views + ?", 1)).Error
}

func (s *BlogService) generateSlug(title string) string {
	slug := strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	var count int64
	s.db.Model(&models.Blog{}).Where("slug LIKE ?", slug+"%").Count(&count)
	if count > 0 {
		slug = fmt.Sprintf("%s-%d", slug, count)
	}

	return slug
}

func (s *BlogService) calculateReadTime(content string) int {
	words := len(strings.Fields(content))
	minutes := int(math.Ceil(float64(words) / 200.0))
	if minutes < 1 {
		minutes = 1
	}
	return minutes
}
