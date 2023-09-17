package server

import (
	"context"
	docs "tokeon-test-task/docs"
	"tokeon-test-task/internal/controllers"
	"tokeon-test-task/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func (s *Server) applyRoutes(
	ctx context.Context,
	mw *middleware.Middleware,
	controllers *controllers.Controllers,
) {
	docs.SwaggerInfo.Host = s.config.ApiAddr

	apiRouter := s.app.Group("/api")

	apiV1Router := apiRouter.Group("/v1")

	apiV1Router.Get("/swagger/*", swagger.HandlerDefault)

	apiV1Router.Get("/health-check", controllers.Common().HealthCheck())
	apiV1Router.Post("/send", controllers.Sender().Send())

	wsRouter := apiV1Router.Group("/ws", mw.Websocket())

	wsRouter.Get("/:id", controllers.Device().Connect(ctx))

	s.app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound) // => 404 "Not Found"
	})
}
