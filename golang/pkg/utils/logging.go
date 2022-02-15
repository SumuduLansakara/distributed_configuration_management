package utils

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogging() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed initializing logger")
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
}
