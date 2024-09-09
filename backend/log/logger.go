package log

import (
	"context"
	"go.uber.org/zap"
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

type logger struct {
	zap *zap.SugaredLogger
}

func NewLogger(appMode string) Logger {
	l, err := zap.NewProduction()
	if appMode == string(DevAppMode) {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}

	return &logger{l.Sugar()}
}

func (logger *logger) Debug(ctx context.Context, msg string) {
	logger.zap.Log(zap.DebugLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}

func (logger *logger) Info(ctx context.Context, msg string) {
	logger.zap.Log(zap.InfoLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}

func (logger *logger) Error(ctx context.Context, msg string) {
	logger.zap.Log(zap.ErrorLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}

func (logger *logger) Fatal(ctx context.Context, msg string) {
	logger.zap.Log(zap.FatalLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}

func (logger *logger) Warn(ctx context.Context, msg string) {
	logger.zap.Log(zap.WarnLevel, msg, GetFields(ctx).String())
	_ = logger.zap.Sync()
}
