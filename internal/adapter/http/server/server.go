package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/ibldzn/spinner-hut/internal/adapter/http/handler"
)

type Config struct {
	Address         string
	Handler         *handler.Handler
	ShutdownTimeout time.Duration
}

type Server struct {
	cfg Config
	srv *http.Server
}

func NewServer(cfg Config) (*Server, error) {
	if cfg.Address == "" {
		return nil, errors.New("server address is required")
	}

	if cfg.Handler == nil {
		return nil, errors.New("handler is required")
	}

	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 5 * time.Second
	}

	return &Server{
		cfg: cfg,
	}, nil
}

func (s *Server) ListenAddr() string {
	return s.cfg.Address
}

func (s *Server) Run(h http.Handler) error {
	s.srv = &http.Server{
		Addr:    s.cfg.Address,
		Handler: h,
	}

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown() error {
	if s.srv == nil {
		return errors.New("server is not running")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	if err := s.srv.Shutdown(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return errors.New("server shutdown timed out")
		} else {
			return err
		}
	}

	return nil
}
