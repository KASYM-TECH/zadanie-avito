//nolint:gochecknoglobals
package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Logger interface {
	Info(ctx context.Context, msg string)
	Warn(ctx context.Context, msg string)
	Debug(ctx context.Context, msg string)
	Error(ctx context.Context, msg string)
	Fatal(ctx context.Context, msg string)
}

type AppMode string

var (
	DevAppMode  AppMode = "dev"
	ProdAppMode AppMode = "prod"
)

type SugaredLogger struct {
	zap *zap.SugaredLogger
}

func NewLogger(appMode string) *SugaredLogger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	var zapCore zapcore.Core
	if appMode == "dev" {
		zapCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(zapcore.Lock(os.Stdout)),
			zapcore.DebugLevel,
		)
	} else {
		zapCore = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(zapcore.Lock(os.Stdout)),
			zapcore.InfoLevel,
		)
	}

	zapLog := zap.New(zapCore)

	return &SugaredLogger{zapLog.Sugar()}
}

func (logger *SugaredLogger) Debug(ctx context.Context, msg string) {
	logger.zap.Log(zap.DebugLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}

func (logger *SugaredLogger) Info(ctx context.Context, msg string) {
	logger.zap.Log(zap.InfoLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}

func (logger *SugaredLogger) Error(ctx context.Context, msg string) {
	logger.zap.Log(zap.ErrorLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}

func (logger *SugaredLogger) Fatal(ctx context.Context, msg string) {
	logger.zap.Log(zap.FatalLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}

func (logger *SugaredLogger) Warn(ctx context.Context, msg string) {
	logger.zap.Log(zap.WarnLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}
