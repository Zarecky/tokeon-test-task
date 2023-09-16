package server

import (
	docs "tokeon-test-task/docs"
	"tokeon-test-task/internal/controller"
	"tokeon-test-task/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (s *Server) applyRoutes(
	mw *middleware.Middleware,
	commonController controller.CommonController,
) {
	docs.SwaggerInfo.Host = s.config.ApiAddr

	apiRouter := s.app.Group("/api")

	apiV1Router := apiRouter.Group("/v1")

	apiV1Router.Get("/swagger/*", swagger.HandlerDefault)

	apiV1Router.Get("/health-check", commonController.HealthCheck())

	s.app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound) // => 404 "Not Found"
	})
}
