package core_http_server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	core_logger "github.com/c-dvoid/teamwork_api_go/internal/core/logger"
	core_http_config "github.com/c-dvoid/teamwork_api_go/internal/core/transport/http/config"
	"go.uber.org/zap"
)

type Server struct {
	handler http.Handler
	config  core_http_config.Config
	log     *core_logger.Logger
}

func NewServer(config core_http_config.Config, log *core_logger.Logger, handler http.Handler) *Server {
	return &Server{
		handler: handler,
		config:  config,
		log:     log,
	}
}

func (s *Server) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    ":" + s.config.Port,
		Handler: s.handler,
	}

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Info("starting HTTP server", zap.String("port", s.config.Port))

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and serve HTTP: %w", err)
		}
	case <-ctx.Done():
		s.log.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		s.log.Info("HTTP server stopped gracefully")
	}

	return nil
}
