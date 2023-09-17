package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) Logger() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := ctx.Next()

		response := ctx.Response()
		code := response.StatusCode()

		loggerExtendedFields := []any{"status_code", code, "ip", ctx.Get("X-Real-IP", ""), "method", ctx.Method(), "url", ctx.OriginalURL()}

		if code < 400 {
			m.logger.With(loggerExtendedFields...).Info("API request")
		}

		return err
	}
}
