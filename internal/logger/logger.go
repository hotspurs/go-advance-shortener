package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	Sugar *zap.SugaredLogger
}

func New(debug bool) *Logger {
	var logger *zap.Logger
	if debug {
		logger, _ = zap.NewDevelopment()
	} else {
		logger, _ = zap.NewProduction()
	}
	sugar := logger.Sugar()

	return &Logger{
		Sugar: sugar,
	}
}

func (l *Logger) Sync() {
	l.Sugar.Sync()
}
