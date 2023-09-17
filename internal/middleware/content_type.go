package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) ApplicationJsonContentType() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Accepts("application/json") != "" {
			return ctx.Next()
		}

		return fiber.NewError(fiber.StatusBadRequest, "Require content types: application/josn")
	}
}
