package middleware

import (
	"github.com/gofiber/fiber/v2"
	"solution/internal/domain/business"
	"strings"
)

func unauthorized(c *fiber.Ctx, details ...string) error {
	if len(details) > 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"message": "Пользователь не авторизован.",
				"detail":  details[0],
				"status":  fiber.StatusUnauthorized,
			},
		)
	}
	return c.Status(fiber.StatusUnauthorized).JSON(
		fiber.Map{
			"message": "Пользователь не авторизован.",
			"status":  fiber.StatusUnauthorized,
		},
	)
}

func TokenAuth(manager business.TokenManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := c.Get("Authorization")
		s := strings.Split(t, " ")
		if len(s) != 2 || s[0] != "Bearer" {
			return unauthorized(c)
		}
		data, err := manager.ValidateToken(s[1])
		if err != nil {
			return unauthorized(c, err.Message)
		}
		c.Locals("sub", data.Sub.String())
		c.Locals("email", data.Email)
		return c.Next()
	}
}
