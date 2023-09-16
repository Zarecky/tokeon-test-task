package middleware

import (
	"tokeon-test-task/pkg/customerror"

	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) ApplicationJsonContentType() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Accepts("application/json") != "" {
			return ctx.Next()
		}

		return customerror.BadRequestError{Message: "Require content types: application/josn"}
	}
}
