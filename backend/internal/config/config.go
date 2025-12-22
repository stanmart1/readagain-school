package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Payment  PaymentConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret            string
	ExpireHours       int
	RefreshExpireDays int
}

type RedisConfig struct {
	URL string
}

type PaymentConfig struct {
	PaystackSecretKey       string
	PaystackPublicKey       string
	FlutterwaveSecretKey    string
	FlutterwavePublicKey    string
	FlutterwaveEncryptionKey string
}

type EmailConfig struct {
	ResendAPIKey string
	FromEmail    string
	FromName     string
	AppURL       string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	refreshDays, _ := strconv.Atoi(getEnv("REFRESH_TOKEN_EXPIRE_DAYS", "7"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8000"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", ""),
			ExpireHours:       expireHours,
			RefreshExpireDays: refreshDays,
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", ""),
		},
		Payment: PaymentConfig{
			PaystackSecretKey:        getEnv("PAYSTACK_SECRET_KEY", ""),
			PaystackPublicKey:        getEnv("PAYSTACK_PUBLIC_KEY", ""),
			FlutterwaveSecretKey:     getEnv("FLUTTERWAVE_SECRET_KEY", ""),
			FlutterwavePublicKey:     getEnv("FLUTTERWAVE_PUBLIC_KEY", ""),
			FlutterwaveEncryptionKey: getEnv("FLUTTERWAVE_ENCRYPTION_KEY", ""),
		},
		Email: EmailConfig{
			ResendAPIKey: getEnv("RESEND_API_KEY", ""),
			FromEmail:    getEnv("FROM_EMAIL", "noreply@readagain.com"),
			FromName:     getEnv("FROM_NAME", "ReadAgain"),
			AppURL:       getEnv("APP_URL", "http://localhost:3000"),
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
