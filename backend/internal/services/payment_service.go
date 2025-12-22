package services

import (
	"fmt"

	"readagain/internal/utils"
)

type PaymentService struct {
	paystack    *PaystackService
	flutterwave *FlutterwaveService
	bankTransfer *BankTransferService
}

func NewPaymentService(paystackKey, flutterwaveKey string, bankTransfer *BankTransferService) *PaymentService {
	return &PaymentService{
		paystack:     NewPaystackService(paystackKey),
		flutterwave:  NewFlutterwaveService(flutterwaveKey),
		bankTransfer: bankTransfer,
	}
}

type PaymentInitializeRequest struct {
	Provider  string                 `json:"provider" validate:"required,oneof=paystack flutterwave bank_transfer"`
	Email     string                 `json:"email" validate:"required,email"`
	Name      string                 `json:"name"`
	Amount    float64                `json:"amount" validate:"required,gt=0"`
	Reference string                 `json:"reference" validate:"required"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type PaymentInitializeResponse struct {
	Provider         string                    `json:"provider"`
	PaymentURL       string                    `json:"payment_url,omitempty"`
	Reference        string                    `json:"reference"`
	BankInstructions *BankTransferInstructions `json:"bank_instructions,omitempty"`
}

func (s *PaymentService) InitializePayment(req PaymentInitializeRequest) (*PaymentInitializeResponse, error) {
	switch req.Provider {
	case "paystack":
		return s.initializePaystack(req)
	case "flutterwave":
		return s.initializeFlutterwave(req)
	case "bank_transfer":
		return s.initializeBankTransfer(req)
	default:
		return nil, utils.NewBadRequestError("Invalid payment provider")
	}
}

func (s *PaymentService) initializePaystack(req PaymentInitializeRequest) (*PaymentInitializeResponse, error) {
	result, err := s.paystack.InitializePayment(req.Email, req.Amount, req.Reference, req.Metadata)
	if err != nil {
		return nil, err
	}

	return &PaymentInitializeResponse{
		Provider:   "paystack",
		PaymentURL: result.Data.AuthorizationURL,
		Reference:  result.Data.Reference,
	}, nil
}

func (s *PaymentService) initializeFlutterwave(req PaymentInitializeRequest) (*PaymentInitializeResponse, error) {
	name := req.Name
	if name == "" {
		name = req.Email
	}

	result, err := s.flutterwave.InitializePayment(req.Email, name, req.Amount, req.Reference, req.Metadata)
	if err != nil {
		return nil, err
	}

	return &PaymentInitializeResponse{
		Provider:   "flutterwave",
		PaymentURL: result.Data.Link,
		Reference:  req.Reference,
	}, nil
}

func (s *PaymentService) initializeBankTransfer(req PaymentInitializeRequest) (*PaymentInitializeResponse, error) {
	instructions := s.bankTransfer.GenerateInstructions(req.Reference, req.Amount)

	return &PaymentInitializeResponse{
		Provider:         "bank_transfer",
		Reference:        req.Reference,
		BankInstructions: instructions,
	}, nil
}

type PaymentVerifyResponse struct {
	Provider      string  `json:"provider"`
	Reference     string  `json:"reference"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
	PaidAt        string  `json:"paid_at"`
	TransactionID string  `json:"transaction_id"`
}

func (s *PaymentService) VerifyPayment(provider, reference string) (*PaymentVerifyResponse, error) {
	switch provider {
	case "paystack":
		return s.verifyPaystack(reference)
	case "flutterwave":
		return s.verifyFlutterwave(reference)
	case "bank_transfer":
		return nil, utils.NewBadRequestError("Bank transfer requires manual verification")
	default:
		return nil, utils.NewBadRequestError("Invalid payment provider")
	}
}

func (s *PaymentService) verifyPaystack(reference string) (*PaymentVerifyResponse, error) {
	result, err := s.paystack.VerifyPayment(reference)
	if err != nil {
		return nil, err
	}

	return &PaymentVerifyResponse{
		Provider:      "paystack",
		Reference:     result.Data.Reference,
		Status:        result.Data.Status,
		Amount:        result.Data.Amount / 100,
		PaidAt:        result.Data.PaidAt,
		TransactionID: fmt.Sprintf("%d", result.Data.ID),
	}, nil
}

func (s *PaymentService) verifyFlutterwave(reference string) (*PaymentVerifyResponse, error) {
	result, err := s.flutterwave.VerifyPayment(reference)
	if err != nil {
		return nil, err
	}

	return &PaymentVerifyResponse{
		Provider:      "flutterwave",
		Reference:     result.Data.TxRef,
		Status:        result.Data.Status,
		Amount:        result.Data.Amount,
		PaidAt:        result.Data.CreatedAt,
		TransactionID: fmt.Sprintf("%d", result.Data.ID),
	}, nil
}
