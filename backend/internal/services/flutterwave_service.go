package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"readagain/internal/utils"
)

type FlutterwaveService struct {
	secretKey string
	baseURL   string
}

func NewFlutterwaveService(secretKey string) *FlutterwaveService {
	return &FlutterwaveService{
		secretKey: secretKey,
		baseURL:   "https://api.flutterwave.com/v3",
	}
}

type FlutterwaveInitializeRequest struct {
	TxRef          string                 `json:"tx_ref"`
	Amount         float64                `json:"amount"`
	Currency       string                 `json:"currency"`
	RedirectURL    string                 `json:"redirect_url"`
	PaymentOptions string                 `json:"payment_options"`
	Customer       FlutterwaveCustomer    `json:"customer"`
	Customizations FlutterwaveCustomization `json:"customizations"`
	Meta           map[string]interface{} `json:"meta,omitempty"`
}

type FlutterwaveCustomer struct {
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber,omitempty"`
	Name        string `json:"name"`
}

type FlutterwaveCustomization struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Logo        string `json:"logo,omitempty"`
}

type FlutterwaveInitializeResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Link string `json:"link"`
	} `json:"data"`
}

type FlutterwaveVerifyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              int64   `json:"id"`
		TxRef           string  `json:"tx_ref"`
		FlwRef          string  `json:"flw_ref"`
		Amount          float64 `json:"amount"`
		Currency        string  `json:"currency"`
		ChargedAmount   float64 `json:"charged_amount"`
		Status          string  `json:"status"`
		PaymentType     string  `json:"payment_type"`
		CreatedAt       string  `json:"created_at"`
		Customer        FlutterwaveCustomer `json:"customer"`
	} `json:"data"`
}

func (s *FlutterwaveService) InitializePayment(email, name string, amount float64, reference string, metadata map[string]interface{}) (*FlutterwaveInitializeResponse, error) {
	payload := FlutterwaveInitializeRequest{
		TxRef:          reference,
		Amount:         amount,
		Currency:       "NGN",
		PaymentOptions: "card,banktransfer,ussd",
		Customer: FlutterwaveCustomer{
			Email: email,
			Name:  name,
		},
		Customizations: FlutterwaveCustomization{
			Title:       "ReadAgain",
			Description: "Book Purchase",
		},
		Meta: metadata,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to marshal request", err)
	}

	req, err := http.NewRequest("POST", s.baseURL+"/payments", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.secretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to initialize payment", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to read response", err)
	}

	var result FlutterwaveInitializeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, utils.NewInternalServerError("Failed to parse response", err)
	}

	if result.Status != "success" {
		return nil, utils.NewBadRequestError(result.Message)
	}

	return &result, nil
}

func (s *FlutterwaveService) VerifyPayment(transactionID string) (*FlutterwaveVerifyResponse, error) {
	url := fmt.Sprintf("%s/transactions/%s/verify", s.baseURL, transactionID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to create request", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.secretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to verify payment", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to read response", err)
	}

	var result FlutterwaveVerifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, utils.NewInternalServerError("Failed to parse response", err)
	}

	if result.Status != "success" {
		return nil, utils.NewBadRequestError(result.Message)
	}

	return &result, nil
}

func (s *FlutterwaveService) ValidateWebhook(signature string, payload []byte) bool {
	return true
}
