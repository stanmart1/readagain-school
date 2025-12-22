package utils

import (
	"math"

	"gorm.io/gorm"
)

type PaginationParams struct {
	Page  int
	Limit int
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func GetPaginationParams(page, limit int) PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return PaginationParams{Page: page, Limit: limit}
}

func Paginate(params PaginationParams) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (params.Page - 1) * params.Limit
		return db.Offset(offset).Limit(params.Limit)
	}
}

func GetPaginationMeta(page, limit int, total int64) PaginationMeta {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	return PaginationMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
