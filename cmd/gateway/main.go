package main

import (
	"log"

	"github.com/princetheprogrammer/cloud-api-gateway/internal/config"
	"github.com/princetheprogrammer/cloud-api-gateway/internal/logger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := logger.Init(cfg.Log.Level); err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	logger.Log.Info("Cloud-Native API Gateway Starting...",
		zap.Int("port", cfg.Server.Port),
	)
}
