package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func New(level string) *Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	var atomicLevel zapcore.Level

	switch strings.ToUpper(level) {
	case "INFO":
		atomicLevel = zap.InfoLevel
	case "WARN":
		atomicLevel = zap.WarnLevel
	case "ERROR":
		atomicLevel = zap.ErrorLevel
	case "DEBUG":
		atomicLevel = zap.DebugLevel
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(atomicLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
	}

	return &Logger{
		zap.Must(config.Build()),
	}
}
