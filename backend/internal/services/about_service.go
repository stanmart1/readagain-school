package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
)

type AboutService struct {
	db *gorm.DB
}

func NewAboutService(db *gorm.DB) *AboutService {
	return &AboutService{db: db}
}

func (s *AboutService) Get() (*models.AboutPage, error) {
	var about models.AboutPage
	if err := s.db.First(&about).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			about = models.AboutPage{
				Title:   "About ReadAgain",
				Content: "Welcome to ReadAgain - Your digital reading platform.",
			}
			s.db.Create(&about)
			return &about, nil
		}
		return nil, err
	}
	return &about, nil
}

func (s *AboutService) Update(updates map[string]interface{}) error {
	var about models.AboutPage
	if err := s.db.First(&about).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			about = models.AboutPage{}
			if err := s.db.Create(&about).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return s.db.Model(&about).Updates(updates).Error
}
