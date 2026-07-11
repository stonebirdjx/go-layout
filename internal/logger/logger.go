package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/stonebirdjx/go-layout/internal/config"
	"github.com/stonebirdjx/go-layout/pkg/consts"
)

// Init initializes the global zap logger based on configuration.
// When FilePath is set, logs are written to both stdout and a rotating file.
// When FilePath is empty, logs are written to stdout only.
func Init(cfg config.LogOptions) *zap.Logger {
	var level zapcore.Level
	switch cfg.Level {
	case consts.LoggerLevelDebug:
		level = zapcore.DebugLevel
	case consts.LoggerLevelInfo:
		level = zapcore.InfoLevel
	case consts.LoggerLevelWarn:
		level = zapcore.WarnLevel
	case consts.LoggerLevelError:
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	var encoder zapcore.Encoder
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Format == consts.LoggerFormatJSON {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var syncers []zapcore.WriteSyncer

	if cfg.EnableStdout {
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}

	if cfg.FilePath != "" {
		fileSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   cfg.FilePath,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		})
		syncers = append(syncers, fileSyncer)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(syncers...),
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger)

	return logger
}
