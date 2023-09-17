package main

import (
	"context"
	"tokeon-test-task/internal/config"
	"tokeon-test-task/internal/server"
	"tokeon-test-task/pkg/initialconfig"
	"tokeon-test-task/pkg/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Loading service config
	cfg := new(config.Config)
	initialconfig.LoadConfig(log.New(), cfg)

	logLevel := log.INFO
	if cfg.EnvCI == "local" {
		logLevel = log.DEBUG
	}

	// Init logger
	logger := log.New(log.WithLogLevel(logLevel))
	defer logger.Sync()

	// Init Server
	srv, err := server.New(logger, cfg)
	if err != nil {
		logger.Fatalf("init server error: %v, ", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Start server
	go func(ctx context.Context) {
		if err := srv.Start(ctx); err != nil {
			logger.Fatalf("start server error: %v, ", err)
		}
	}(ctx)

	// Wait system signals
	<-sig

	// Stop server
	cancel()
	srv.Stop()
}
