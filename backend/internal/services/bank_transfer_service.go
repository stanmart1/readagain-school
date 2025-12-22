package services

import (
	"fmt"
	"time"

	"gorm.io/gorm"

	"readagain/internal/models"
	"readagain/internal/utils"
)

type BankTransferService struct {
	db *gorm.DB
}

func NewBankTransferService(db *gorm.DB) *BankTransferService {
	return &BankTransferService{db: db}
}

type BankAccount struct {
	BankName      string `json:"bank_name"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
}

type BankTransferInstructions struct {
	Reference    string        `json:"reference"`
	Amount       float64       `json:"amount"`
	BankAccounts []BankAccount `json:"bank_accounts"`
	Instructions string        `json:"instructions"`
	ExpiresAt    time.Time     `json:"expires_at"`
}

func (s *BankTransferService) GetBankAccounts() []BankAccount {
	return []BankAccount{
		{
			BankName:      "GTBank",
			AccountNumber: "0123456789",
			AccountName:   "ReadAgain Limited",
		},
		{
			BankName:      "Access Bank",
			AccountNumber: "0987654321",
			AccountName:   "ReadAgain Limited",
		},
	}
}

func (s *BankTransferService) GenerateInstructions(reference string, amount float64) *BankTransferInstructions {
	expiresAt := time.Now().Add(24 * time.Hour)

	instructions := `
Please transfer the exact amount to any of the bank accounts listed above.
Use the reference code as your payment description/narration.
After payment, upload your payment proof for verification.
Payment will be confirmed within 1-24 hours.
	`

	return &BankTransferInstructions{
		Reference:    reference,
		Amount:       amount,
		BankAccounts: s.GetBankAccounts(),
		Instructions: instructions,
		ExpiresAt:    expiresAt,
	}
}

func (s *BankTransferService) UploadProof(orderID uint, proofURL string) error {
	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return utils.NewNotFoundError("Order not found")
	}

	if order.PaymentMethod != "bank_transfer" {
		return utils.NewBadRequestError("Order payment method is not bank transfer")
	}

	order.Notes = fmt.Sprintf("%s | Proof: %s | Status: pending_verification", order.Notes, proofURL)

	if err := s.db.Save(&order).Error; err != nil {
		return utils.NewInternalServerError("Failed to update order", err)
	}

	return nil
}

func (s *BankTransferService) VerifyPayment(orderID uint, verified bool) error {
	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return utils.NewNotFoundError("Order not found")
	}

	if order.PaymentMethod != "bank_transfer" {
		return utils.NewBadRequestError("Order payment method is not bank transfer")
	}

	if verified {
		order.Status = "completed"
		order.Notes = fmt.Sprintf("%s | Verified", order.Notes)
	} else {
		order.Status = "failed"
		order.Notes = fmt.Sprintf("%s | Verification failed", order.Notes)
	}

	if err := s.db.Save(&order).Error; err != nil {
		return utils.NewInternalServerError("Failed to update order", err)
	}

	return nil
}

func (s *BankTransferService) GetPendingVerifications(page, limit int) ([]models.Order, *utils.PaginationMeta, error) {
	params := utils.GetPaginationParams(page, limit)

	query := s.db.Model(&models.Order{}).
		Where("payment_method = ? AND notes LIKE ?", "bank_transfer", "%pending_verification%").
		Preload("User").
		Preload("Items.Book")

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to count orders", err)
	}

	var orders []models.Order
	if err := query.Scopes(utils.Paginate(params)).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, nil, utils.NewInternalServerError("Failed to fetch orders", err)
	}

	meta := utils.GetPaginationMeta(params.Page, params.Limit, total)
	return orders, &meta, nil
}
