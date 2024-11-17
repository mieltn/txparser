package logger

import (
	"fmt"
	"log/slog"
	"os"
)

const (
	devMode  = "dev"
	prodMode = "prod"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type logger struct {
	l *slog.Logger
}

func New(mode string) Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})
	switch mode {
	case prodMode:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	l := slog.New(handler)
	return &logger{
		l: l,
	}
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.l.Info(fmt.Sprintf(format, args...))
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.l.Debug(fmt.Sprintf(format, args...))
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.l.Error(fmt.Sprintf(format, args...))
}
