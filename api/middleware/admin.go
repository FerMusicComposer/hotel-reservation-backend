package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AdminAuth(c *fiber.Ctx) error {
	userClaims := c.Locals("userClaims").(jwt.MapClaims)

	if role, ok := userClaims["role"].(string); ok {
		if role == "admin" {
			return c.Next()
		}
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "unauthorized",
	})
}
