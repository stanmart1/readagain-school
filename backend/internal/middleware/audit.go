package middleware

import (
	"github.com/gofiber/fiber/v2"

	"readagain/internal/services"
)

func AuditMiddleware(auditService *services.AuditService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Locals("auditService", auditService)
		return c.Next()
	}
}

func LogAudit(c *fiber.Ctx, action, entityType string, entityID uint, oldValue, newValue string) {
	auditService, ok := c.Locals("auditService").(*services.AuditService)
	if !ok {
		return
	}

	userID, ok := c.Locals("userID").(uint)
	if !ok {
		userID = 0
	}

	ipAddress := c.IP()
	userAgent := c.Get("User-Agent")

	go auditService.Log(userID, action, entityType, entityID, oldValue, newValue, ipAddress, userAgent)
}
