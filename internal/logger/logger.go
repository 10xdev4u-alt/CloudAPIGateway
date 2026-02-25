package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(level string) error {
	var config zap.Config
	if level == "debug" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Log = logger
	return nil
}
