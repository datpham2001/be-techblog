package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	pkgConfig "github.com/datpham2001/techblog/internal/config"
	pkgLogger "github.com/datpham2001/techblog/internal/logger"
	"github.com/datpham2001/techblog/internal/server"
)

type ServerManager struct {
	HTTPServer *server.Server
}

var (
	cfg    *pkgConfig.Config = &pkgConfig.Config{}
	logger *pkgLogger.Logger
)

func init() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	configPath := filepath.Join(workingDir, "configs")
	if err := pkgConfig.Load(configPath, cfg); err != nil {
		log.Fatalf("failed to load config from path %s: %v", configPath, err)
	}

	logger = pkgLogger.Initalize(cfg)
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverManager := &ServerManager{}
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := serverManager.shutdownHTTPServer(shutdownCtx); err != nil {
			logger.Fatal(err)
		}

		// close database connection
		// close redis connection

		logger.Info("Server shutdown gracefully")
	}()

	go func() {
		logger.Info("Starting Techblog server...")
		if err := serverManager.startHTTPServer(ctx); err != nil {
			logger.Fatalf("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
}

func (s *ServerManager) startHTTPServer(ctx context.Context) error {
	s.HTTPServer = server.New(cfg, logger)
	return s.HTTPServer.Start(ctx)
}

func (s *ServerManager) shutdownHTTPServer(ctx context.Context) error {
	if s.HTTPServer == nil {
		return nil
	}

	return s.HTTPServer.Shutdown(ctx)
}
