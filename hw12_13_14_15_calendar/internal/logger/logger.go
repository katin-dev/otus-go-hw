package logger

import (
	"fmt"

	"go.uber.org/zap"
)

const (
	ENV_PROD = "prod"
	ENV_DEV  = "dev"
)

type Logger struct {
	file string
	logg *zap.SugaredLogger
}

func New(file, env string) (*Logger, error) {
	var cfg zap.Config

	switch env {
	case ENV_PROD:
		cfg = zap.NewProductionConfig()
	case ENV_DEV:
		cfg = zap.NewDevelopmentConfig()
	default:
		return nil, fmt.Errorf("env '%s' must be one of [dev, prod]", env)
	}

	cfg.OutputPaths = []string{file}

	zapLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger := &Logger{
		file: file,
		logg: zapLogger.Sugar(),
	}

	return logger, nil
}

func (l *Logger) Debug(msg string, params ...interface{}) {
	l.logg.Debugw(msg, params...)
}

func (l *Logger) Info(msg string, params ...interface{}) {
	l.logg.Infow(msg, params...)
}

func (l *Logger) Warn(msg string, params ...interface{}) {
	l.logg.Warnw(msg, params...)
}

func (l *Logger) Error(msg string, params ...interface{}) {
	l.logg.Errorw(msg, params...)
}
