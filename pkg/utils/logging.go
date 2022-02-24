package utils

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogging() {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	logger, err := config.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed initializing logger")
		os.Exit(1)
	}
	zap.ReplaceGlobals(logger)
}
