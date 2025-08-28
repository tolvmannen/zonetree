package logger

import (
	"log/slog"
	"os"
)

type SlogLogger struct {
	l *slog.Logger
}

func NewSlogLogger(l *slog.Logger) SlogLogger {
	return SlogLogger{l: l}
}

func PrintDebugLog() SlogLogger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	//s := slog.New(slog.NewTextHandler(os.Stdout, opts))
	s := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	return SlogLogger{l: s}
}

func (s SlogLogger) Debug(msg string, args ...any) {
	s.l.Debug(msg, args...)
}

func (s SlogLogger) Info(msg string, args ...any) {
	s.l.Info(msg, args...)
}

func (s SlogLogger) Warn(msg string, args ...any) {
	s.l.Warn(msg, args...)
}

func (s SlogLogger) Error(msg string, args ...any) {
	s.l.Error(msg, args...)
}

func (s SlogLogger) With(args ...any) Logger {
	return SlogLogger{l: s.l.With(args...)}
}
