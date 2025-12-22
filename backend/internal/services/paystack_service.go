package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"readagain/internal/utils"
)

type PaystackService struct {
	secretKey string
	baseURL   string
}

func NewPaystackService(secretKey string) *PaystackService {
	return &PaystackService{
		secretKey: secretKey,
		baseURL:   "https://api.paystack.co",
	}
}

type PaystackInitializeRequest struct {
	Email     string  `json:"email"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
	CallbackURL string `json:"callback_url,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

type PaystackInitializeResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type PaystackVerifyResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		ID              int64   `json:"id"`
		Status          string  `json:"status"`
		Reference       string  `json:"reference"`
		Amount          float64 `json:"amount"`
		PaidAt          string  `json:"paid_at"`
		Channel         string  `json:"channel"`
		Currency        string  `json:"currency"`
		TransactionDate string  `json:"transaction_date"`
		Customer        struct {
			Email string `json:"email"`
		} `json:"customer"`
	} `json:"data"`
}

func (s *PaystackService) InitializePayment(email string, amount float64, reference string, metadata map[string]interface{}) (*PaystackInitializeResponse, error) {
	amountInKobo := amount * 100

	payload := PaystackInitializeRequest{
		Email:     email,
		Amount:    amountInKobo,
		Reference: reference,
		Metadata:  metadata,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, utils.NewInternalServerError("Failed to marshal request", err)
	}

	req, err := http.NewRequest("POST", s.baseURL+"/transaction/initialize", bytes.NewBuffer(jsonData))
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

	var result PaystackInitializeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, utils.NewInternalServerError("Failed to parse response", err)
	}

	if !result.Status {
		return nil, utils.NewBadRequestError(result.Message)
	}

	return &result, nil
}

func (s *PaystackService) VerifyPayment(reference string) (*PaystackVerifyResponse, error) {
	url := fmt.Sprintf("%s/transaction/verify/%s", s.baseURL, reference)

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

	var result PaystackVerifyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, utils.NewInternalServerError("Failed to parse response", err)
	}

	if !result.Status {
		return nil, utils.NewBadRequestError(result.Message)
	}

	return &result, nil
}

func (s *PaystackService) ValidateWebhook(signature string, payload []byte) bool {
	return true
}
