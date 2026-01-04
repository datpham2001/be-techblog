package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/datpham2001/techblog/internal/config"
	"github.com/datpham2001/techblog/internal/logger"
	"github.com/datpham2001/techblog/internal/middlewares"
	"github.com/gin-gonic/gin"
)

const (
	SHUTDOWN_TIMEOUT = 5 * time.Second
)

type Server struct {
	router *gin.Engine
	cfg    *config.Config
	server *http.Server
	logger *logger.Logger
}

func New(cfg *config.Config, logger *logger.Logger) *Server {
	if cfg.Server.Env == "local" || cfg.Server.Env == "development" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	setupMiddlewares(router, cfg, logger)

	return &Server{
		router: router,
		cfg:    cfg,
		logger: logger,
	}
}

func setupMiddlewares(router *gin.Engine, cfg *config.Config, l *logger.Logger) {
	// Recovery middleware (must be first)
	router.Use(gin.Recovery())

	// Logger middleware (using logrus)
	router.Use(middlewares.Logger(l))

	// Add other middlewares here (CORS, rate limit, auth, etc.)
}

func (s *Server) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	s.server = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port),
		Handler: s.router,
	}
	s.server.RegisterOnShutdown(cancel)

	go func() {
		addr := fmt.Sprintf("%s:%s", s.cfg.Server.Host, s.cfg.Server.Port)
		s.logger.Infof("Server is listening on %s", addr)

		var err error
		if s.cfg.Server.TLS.Enable {
			err = s.server.ListenAndServeTLS(s.cfg.Server.TLS.CertFile, s.cfg.Server.TLS.KeyFile)
		} else {
			err = s.server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("failed to start http server: %v", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		if closeErr := s.server.Close(); closeErr != nil {
			return fmt.Errorf("could not stop server gracefully and force to close: %w", closeErr)
		}

		return fmt.Errorf("could not stop server gracefully: %w", err)
	}

	return nil
}
