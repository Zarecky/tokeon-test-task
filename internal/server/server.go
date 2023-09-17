package server

import (
	"context"
	"fmt"
	"strings"
	"time"

	"tokeon-test-task/internal/config"
	"tokeon-test-task/internal/controllers"
	"tokeon-test-task/internal/middleware"
	"tokeon-test-task/internal/services"
	"tokeon-test-task/pkg/hc"
	"tokeon-test-task/pkg/log"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	config *config.Config
	logger log.Logger

	hc *hc.Server

	app *fiber.App

	// Dependencies
	services *services.Services
}

func New(logger log.Logger, cfg *config.Config) (*Server, error) {
	s := &Server{
		logger: logger,
		config: cfg,
	}

	return s, nil
}

func (s *Server) Start(ctx context.Context) error {
	defer s.Stop()

	// Init internal services
	if err := s.initInternalServices(ctx); err != nil {
		return fmt.Errorf("failed to init internal services: %w", err)
	}

	if s.config.EnvCI != "local" {
		go func() {
			s.startHealthCheckServer()

			for {
				time.Sleep(time.Minute * 5)
				s.logger.Info("Interval service logging")
			}
		}()
	}

	s.logger.Info("server started successfully")

	// Wait context
	<-ctx.Done()

	return nil
}

func (s *Server) Stop() {
	// stop hc
	if s.hc != nil {
		s.hc.Stop(context.Background())
	}

	s.logger.Info("server stopped")
}

func (s *Server) startHealthCheckServer() {
	// Init HC Server
	s.hc = hc.NewServer(s.logger, s.config.HealthCheck)

	// Register services
	s.hc.RegisterService(s.config.ServiceName, hc.NewService(0, nil, nil))

	// Start HC Server
	go s.hc.Start()
}

func (s *Server) initInternalServices(ctx context.Context) error {
	// Init services
	var err error
	s.services, err = services.New(s.logger, s.config)
	if err != nil {
		return fmt.Errorf("failed to init services: %w", err)
	}

	// init middleware
	mw := middleware.New(s.logger, s.config)

	// Create http server
	s.app = fiber.New(fiber.Config{
		EnableSplittingOnParsers:     true,
		DisableStartupMessage:        true,
		Immutable:                    true,
		ErrorHandler:                 mw.ErrorHandler(),
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		BodyLimit:                    5 * 1024 * 1024 * 1024,
	})

	s.app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
	}))
	s.app.Use(mw.Logger())

	if s.config.EnvCI == "local" {
		s.app.Use(cors.New(cors.Config{
			Next:         nil,
			AllowOrigins: "*",
			AllowMethods: strings.Join([]string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodHead,
				fiber.MethodPut,
				fiber.MethodDelete,
				fiber.MethodPatch,
			}, ","),
			AllowHeaders:     "",
			AllowCredentials: false,
			ExposeHeaders:    "",
			MaxAge:           0,
		}))
	}

	validator := validator.New()

	// init and apply controllers
	controllers := controllers.New(s.logger, validator, s.services.Device(), s.services.Device())

	s.applyRoutes(
		ctx,
		mw,
		controllers,
	)

	// start rest api server
	go func() {
		if err := s.app.Listen(fmt.Sprintf(":%d", s.config.Port)); err != nil {
			s.logger.Fatalf("failed to start http server: %v", err)
		}
	}()

	return nil
}
