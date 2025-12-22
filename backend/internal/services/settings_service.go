package services

import (
	"gorm.io/gorm"

	"readagain/internal/models"
)

type SettingsService struct {
	db *gorm.DB
}

func NewSettingsService(db *gorm.DB) *SettingsService {
	return &SettingsService{db: db}
}

func (s *SettingsService) GetByCategory(category string) ([]models.SystemSettings, error) {
	var settings []models.SystemSettings
	query := s.db.Model(&models.SystemSettings{})

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Order("key ASC").Find(&settings).Error; err != nil {
		return nil, err
	}

	return settings, nil
}

func (s *SettingsService) GetByKey(key string) (*models.SystemSettings, error) {
	var setting models.SystemSettings
	if err := s.db.Where("key = ?", key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (s *SettingsService) Set(key, value, category, description string) error {
	var setting models.SystemSettings
	err := s.db.Where("key = ?", key).First(&setting).Error

	if err == gorm.ErrRecordNotFound {
		setting = models.SystemSettings{
			Key:         key,
			Value:       value,
			Category:    category,
			Description: description,
		}
		return s.db.Create(&setting).Error
	}

	if err != nil {
		return err
	}

	return s.db.Model(&setting).Updates(map[string]interface{}{
		"value":       value,
		"category":    category,
		"description": description,
	}).Error
}

func (s *SettingsService) Delete(key string) error {
	return s.db.Where("key = ?", key).Delete(&models.SystemSettings{}).Error
}

func (s *SettingsService) GetEmailSettings() (map[string]string, error) {
	settings, err := s.GetByCategory("email")
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	return result, nil
}

func (s *SettingsService) GetPaymentSettings() (map[string]string, error) {
	settings, err := s.GetByCategory("payment")
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}

	return result, nil
}

func (s *SettingsService) UpdateEmailSettings(resendAPIKey, fromEmail, fromName string) error {
	if err := s.Set("email.resend_api_key", resendAPIKey, "email", "Resend API Key"); err != nil {
		return err
	}
	if err := s.Set("email.from_email", fromEmail, "email", "From Email Address"); err != nil {
		return err
	}
	if err := s.Set("email.from_name", fromName, "email", "From Name"); err != nil {
		return err
	}
	return nil
}

func (s *SettingsService) UpdatePaymentSettings(paystackKey, flutterwaveKey, bankName, accountNumber, accountName string) error {
	if err := s.Set("payment.paystack_secret_key", paystackKey, "payment", "Paystack Secret Key"); err != nil {
		return err
	}
	if err := s.Set("payment.flutterwave_secret_key", flutterwaveKey, "payment", "Flutterwave Secret Key"); err != nil {
		return err
	}
	if err := s.Set("payment.bank_name", bankName, "payment", "Bank Name for Transfers"); err != nil {
		return err
	}
	if err := s.Set("payment.account_number", accountNumber, "payment", "Bank Account Number"); err != nil {
		return err
	}
	if err := s.Set("payment.account_name", accountName, "payment", "Bank Account Name"); err != nil {
		return err
	}
	return nil
}
