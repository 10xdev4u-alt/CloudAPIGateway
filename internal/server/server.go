package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/princetheprogrammer/cloud-api-gateway/internal/config"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/logger"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/middleware"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/proxy"
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
	s.router.AddRoute(http.MethodGet, "/health", "self", s.healthHandler)

	for _, rConfig := range s.cfg.Routes {
		p, err := proxy.NewReverseProxy(rConfig.Target)
		if err != nil {
			logger.Log.Error("Failed to create proxy for route", zap.String("path", rConfig.Path), zap.Error(err))
			continue
		}

		var handler http.Handler = p
		if rConfig.StripPrefix {
			handler = http.StripPrefix(rConfig.Path, p)
		}

		s.router.AddRoute(rConfig.Method, rConfig.Path, rConfig.Target, handler.ServeHTTP)
		logger.Log.Info("Registered route", zap.String("method", rConfig.Method), zap.String("path", rConfig.Path), zap.String("target", rConfig.Target))
	}

	chain := middleware.NewChain(middleware.Logging())
	
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Server.Port),
		Handler:      chain.Then(s),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Log.Info("Starting HTTP Server", zap.Int("port", s.cfg.Server.Port))
	return s.httpServer.ListenAndServe()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route, found := s.router.Match(r.Method, r.URL.Path)
	if !found {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Route not found: %s %s\n", r.Method, r.URL.Path)
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
