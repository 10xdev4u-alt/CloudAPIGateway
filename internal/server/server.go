package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/princetheprogrammer/cloud-api-gateway/internal/config"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/logger"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handler)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Log.Info("Starting HTTP Server", zap.Int("port", s.cfg.Server.Port))
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Log.Info("Shutting down HTTP Server...")
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Cloud-Native API Gateway: Request received on %s
", r.URL.Path)
}
