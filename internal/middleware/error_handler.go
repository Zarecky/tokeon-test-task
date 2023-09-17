package middleware

import (
	"fmt"
	"slices"

	"tokeon-test-task/internal/errors"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func (m *Middleware) ErrorHandler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError
		response := ErrorResponse{
			Error: "Internal server error",
		}

		e0, ok := err.(*fiber.Error)
		if ok {
			code = e0.Code
			response = ErrorResponse{
				Error: e0.Message,
			}
		}

		badRequestErrors := []error{
			errors.ErrDeviceAlreadyRegistered,
			errors.ErrDeviceNotFound,
		}

		if slices.Contains(badRequestErrors, err) {
			code = fiber.StatusBadRequest
			response = ErrorResponse{
				Error: err.Error(),
			}
		}

		loggerExtendedFields := []any{"status_code", code, "ip", ctx.Get("X-Real-IP", ""), "method", ctx.Method(), "url", ctx.OriginalURL()}

		errText := fmt.Sprintf("%+v", err)

		switch {
		case code >= 500:
			m.logger.With(loggerExtendedFields...).Error(errText)
		case code >= 400:
			m.logger.With(loggerExtendedFields...).Warn(errText)
		}

		// Set Content-Type: text/plain; charset=utf-8
		ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

		// Return statuscode with error message
		return ctx.Status(code).JSON(response)
	}
}
