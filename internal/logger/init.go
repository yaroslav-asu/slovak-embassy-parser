package logger

import (
	"go.uber.org/zap"
	"log"
	"main/internal/utils/vars"
)

func Init() {
	var logger *zap.Logger
	var err error
	switch vars.RunningMode {
	case "dev":
		logger, err = zap.NewDevelopment()
	case "prod":
		logger, err = zap.NewProduction()
	default:
		logger, err = zap.NewProduction()
	}

	if err != nil {
		log.Fatal("Failed to initialize logger")
	}
	defer logger.Sync()
	zap.ReplaceGlobals(logger)
}
