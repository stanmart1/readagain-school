package middleware

import (
	"github.com/gofiber/fiber/v2"

	"readagain/internal/utils"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if appErr, ok := err.(*utils.AppError); ok {
		code = appErr.Code
		message = appErr.Message
		if appErr.Err != nil {
			utils.ErrorLogger.Printf("Error: %v", appErr.Err)
		}
	} else if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	} else {
		utils.ErrorLogger.Printf("Unexpected error: %v", err)
	}

	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}
