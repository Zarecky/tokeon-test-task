package middleware

import (
	"errors"
	"fmt"
	"tokeon-test-task/internal/repos/common"

	"tokeon-test-task/pkg/customerror"

	"tokeon-test-task/pkg/log"
	"tokeon-test-task/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func (m *Middleware) ErrorHandler(logger log.Logger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		// Status code defaults to 500
		code := fiber.StatusInternalServerError
		response := ErrorResponse{
			Error: "Internal server error",
		}

		e1, ok := err.(customerror.ForbiddenError)
		if ok {
			code = fiber.StatusForbidden
			response = ErrorResponse{
				Error: e1.Message,
			}
		}

		e2, ok := err.(customerror.BadRequestError)
		if ok {
			code = fiber.StatusBadRequest
			response = ErrorResponse{
				Error: e2.Message,
			}
		}

		badRequestErrors := []error{
			common.ErrEmptyID,
		}
		existError := utils.Contains(badRequestErrors, err, func(v1, v2 error) bool {
			return errors.Is(v1, v2)
		})

		if existError {
			code = fiber.StatusBadRequest
			response = ErrorResponse{
				Error: err.Error(),
			}
		}

		notFoundErrors := []error{
			common.ErrNotFound,
		}

		existError = utils.Contains(notFoundErrors, err, func(v1, v2 error) bool {
			return errors.Is(v1, v2)
		})

		if existError {
			code = fiber.StatusNotFound
			response = ErrorResponse{
				Error: err.Error(),
			}
		}

		loggerExtendedFields := []any{"status_code", code, "ip", ctx.Get("X-Real-IP", ""), "method", ctx.Method(), "url", ctx.OriginalURL()}

		errText := fmt.Sprintf("%+v", err)

		switch {
		case code >= 500:
			logger.With(loggerExtendedFields...).Error(errText)
		case code >= 400:
			logger.With(loggerExtendedFields...).Warn(errText)
		}

		// Set Content-Type: text/plain; charset=utf-8
		ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)

		// Return statuscode with error message
		return ctx.Status(code).JSON(response)
	}
}
