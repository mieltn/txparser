package server

import (
	"context"
	"github.com/mieltn/txparser/internal/config"
	"github.com/mieltn/txparser/internal/logger"
	"net/http"
	"os"
	"time"
)

type server struct {
	l   logger.Logger
	srv *http.Server
}

func New(l logger.Logger, cfg config.Config, h http.Handler) *server {
	return &server{
		l: l,
		srv: &http.Server{
			Addr:         cfg.App.Port,
			Handler:      h,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
	}
}

func (s *server) Run(ctx context.Context) {
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.l.Errorf("failed to serve: %v", err)
			time.Sleep(time.Second * 1)
			os.Exit(1)
		}
	}()
	s.l.Infof("server started")

	<-ctx.Done()
	s.l.Infof("interrupt signal received")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.srv.Shutdown(ctxShutdown); err != nil {
		s.l.Errorf("failed to shutdown: %v", err)
	}
	s.l.Infof("server stopped")
}
