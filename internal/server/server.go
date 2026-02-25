package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/princetheprogrammer/cloud-api-gateway/internal/config"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/logger"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/router"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	cfg        *config.Config
	router     *router.Router
}

func New(cfg *config.Config) *Server {
	return &Server{
		cfg:    cfg,
		router: router.New(),
	}
}

func (s *Server) Start() error {
	s.router.AddRoute("/health", "self", s.healthHandler)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler:      s,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Log.Info("Starting HTTP Server", zap.Int("port", s.cfg.Server.Port))
	return s.httpServer.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, found := s.router.Match(r.URL.Path)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Route not found: %s\n", r.URL.Path)
		return
	}

	route.Handler(w, r)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Cloud-Native API Gateway: Healthy\n")
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Log.Info("Shutting down HTTP Server...")
	return s.httpServer.Shutdown(ctx)
}
