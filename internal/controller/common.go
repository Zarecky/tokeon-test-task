package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type PageOptionsDto struct {
	Page     *uint64 `validate:"omitempty,gte=1" query:"page" json:"page"`
	PageSize *uint64 `validate:"omitempty,gte=1" query:"page_size" json:"pageSize"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CommonController interface {
	HealthCheck() fiber.Handler
}

type commonController struct{}

func NewCommonController() CommonController {
	return &commonController{}
}

type healthCheckResponse struct {
	Message string `json:"message"`
}

// HealthCkeck godoc
//
//	@Summary		health check
//	@Description	health check
//	@Tags			common
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	healthCheckResponse
//	@Router			/api/v1/health-check [get]
func (c *commonController) HealthCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.SendStatus(http.StatusOK)
		return c.JSON(healthCheckResponse{Message: "Hello World"})
	}
}
