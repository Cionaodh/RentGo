package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing X-User-ID header",
			})
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid user ID format",
			})
		}

		// Сохраняем в контекст Fiber
		c.Locals("userID", userID)

		return c.Next()
	}
}
