package utils

import "fmt"

const (
	CacheTTLShort  = 300   // 5 minutes
	CacheTTLMedium = 1800  // 30 minutes
	CacheTTLLong   = 3600  // 1 hour
	CacheTTLDay    = 86400 // 24 hours
)

func BookListKey(page, limit int, filters string) string {
	return fmt.Sprintf("books:list:%d:%d:%s", page, limit, filters)
}

func BookDetailKey(id uint) string {
	return fmt.Sprintf("books:detail:%d", id)
}

func CategoryListKey() string {
	return "categories:list"
}

func CategoryDetailKey(id uint) string {
	return fmt.Sprintf("categories:detail:%d", id)
}

func AuthorListKey() string {
	return "authors:list"
}

func AuthorDetailKey(id uint) string {
	return fmt.Sprintf("authors:detail:%d", id)
}

func FeaturedBooksKey() string {
	return "books:featured"
}

func BestsellersKey() string {
	return "books:bestsellers"
}

func NewReleasesKey() string {
	return "books:new-releases"
}

func UserPermissionsKey(userID uint) string {
	return fmt.Sprintf("user:permissions:%d", userID)
}

func SystemSettingsKey(category string) string {
	if category == "" {
		return "settings:all"
	}
	return fmt.Sprintf("settings:%s", category)
}

func BlogListKey(page, limit int, status string) string {
	return fmt.Sprintf("blogs:list:%d:%d:%s", page, limit, status)
}

func BlogDetailKey(slug string) string {
	return fmt.Sprintf("blogs:detail:%s", slug)
}

func FAQListKey(category string) string {
	if category == "" {
		return "faqs:list:all"
	}
	return fmt.Sprintf("faqs:list:%s", category)
}

func TestimonialListKey() string {
	return "testimonials:list"
}
