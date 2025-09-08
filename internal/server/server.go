package server

import (
	"context"
	"fmt"
	"net/http"

	"subscriptions-api/cmd/config"
	"subscriptions-api/pkg/utils/logger"
)

type StoppableServer struct {
	Server *http.Server
	Log    *logger.Logger
}

func NewServer(cfg *config.Config, log *logger.Logger, handler http.Handler) *StoppableServer {
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
	return &StoppableServer{Server: srv, Log: log}
}

// Start запускает сервер в отдельной горутине
func (s *StoppableServer) Start() {
	go func() {
		s.Log.Info("Starting server", "addr", s.Server.Addr)
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Log.Error("Server error", "error", err)
		}
	}()
}

// Shutdown корректно завершает работу сервера
func (s *StoppableServer) Shutdown(ctx context.Context) error {
	s.Log.Info("Shutting down server gracefully")
	return s.Server.Shutdown(ctx)
}

// Stop форсированно закрывает сервер
func (s *StoppableServer) Stop() {
	s.Log.Warn("Force stopping server")
	_ = s.Server.Close()
}
