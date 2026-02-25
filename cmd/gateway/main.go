package main

import (

	"context"

	"log"

	"net/http"

	"os"

	"os/signal"

	"syscall"

	"time"



	"github.com/princetheprogrammer/cloud-api-gateway/internal/config"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/logger"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/server"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/wasm"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := logger.Init(cfg.Log.Level); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	// Initialize Wasm Manager
	wasmManager := wasm.NewManager(ctx)
	defer wasmManager.Close(ctx)

	if err := wasmManager.LoadPlugin(ctx, "example", "plugins/example.wasm"); err != nil {
		logger.Log.Error("failed to load example plugin", zap.Error(err))
	} else {
		if err := wasmManager.RunGreet(ctx, "example"); err != nil {
			logger.Log.Error("failed to run greet", zap.Error(err))
		}
	}

	srv := server.New(cfg)

	go func() {
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("failed to start server", zap.Error(err))
		}
	}()

	logger.Log.Info("Cloud-Native API Gateway Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("Shutting down Gateway...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Fatal("Gateway shutdown failed", zap.Error(err))
	}

	logger.Log.Info("Gateway exited properly")
}
