package services

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/resend/resend-go/v2"

	"readagain/internal/utils"
)

type EmailService struct {
	client       *resend.Client
	fromEmail    string
	fromName     string
	appURL       string
	templatesDir string
}

func NewEmailService(apiKey, fromEmail, fromName, appURL string) *EmailService {
	client := resend.NewClient(apiKey)
	return &EmailService{
		client:       client,
		fromEmail:    fromEmail,
		fromName:     fromName,
		appURL:       appURL,
		templatesDir: "./templates/emails",
	}
}

type WelcomeEmailData struct {
	Name   string
	AppURL string
}

type OrderConfirmationData struct {
	Name        string
	OrderNumber string
	Books       []OrderBook
	Total       string
	AppURL      string
}

type OrderBook struct {
	Title  string
	Author string
	Price  string
}

type PasswordResetData struct {
	Name     string
	ResetURL string
}

func (s *EmailService) SendWelcomeEmail(toEmail, name string) error {
	data := WelcomeEmailData{
		Name:   name,
		AppURL: s.appURL,
	}

	html, err := s.renderTemplate("welcome.html", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{toEmail},
		Subject: "Welcome to ReadAgain!",
		Html:    html,
	}

	_, err = s.client.Emails.Send(params)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to send welcome email: %v", err)
		return utils.NewInternalServerError("Failed to send email", err)
	}

	utils.InfoLogger.Printf("Welcome email sent to %s", toEmail)
	return nil
}

func (s *EmailService) SendOrderConfirmation(toEmail, name, orderNumber string, books []OrderBook, total string) error {
	data := OrderConfirmationData{
		Name:        name,
		OrderNumber: orderNumber,
		Books:       books,
		Total:       total,
		AppURL:      s.appURL,
	}

	html, err := s.renderTemplate("order_confirmation.html", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{toEmail},
		Subject: fmt.Sprintf("Order Confirmed - %s", orderNumber),
		Html:    html,
	}

	_, err = s.client.Emails.Send(params)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to send order confirmation email: %v", err)
		return utils.NewInternalServerError("Failed to send email", err)
	}

	utils.InfoLogger.Printf("Order confirmation email sent to %s", toEmail)
	return nil
}

func (s *EmailService) SendPasswordReset(toEmail, name, resetToken string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", s.appURL, resetToken)

	data := PasswordResetData{
		Name:     name,
		ResetURL: resetURL,
	}

	html, err := s.renderTemplate("password_reset.html", data)
	if err != nil {
		return err
	}

	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", s.fromName, s.fromEmail),
		To:      []string{toEmail},
		Subject: "Reset Your Password - ReadAgain",
		Html:    html,
	}

	_, err = s.client.Emails.Send(params)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to send password reset email: %v", err)
		return utils.NewInternalServerError("Failed to send email", err)
	}

	utils.InfoLogger.Printf("Password reset email sent to %s", toEmail)
	return nil
}

func (s *EmailService) renderTemplate(templateName string, data interface{}) (string, error) {
	templatePath := filepath.Join(s.templatesDir, templateName)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to parse template %s: %v", templateName, err)
		return "", utils.NewInternalServerError("Failed to parse email template", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		utils.ErrorLogger.Printf("Failed to execute template %s: %v", templateName, err)
		return "", utils.NewInternalServerError("Failed to render email template", err)
	}

	return buf.String(), nil
}
