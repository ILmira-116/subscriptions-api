package logger

import (
	"log/slog"
	"os"
	"subscriptions-api/cmd/config"
)

// Logger — обёртка над slog.Logger
type Logger struct {
	log *slog.Logger
}

// New создаёт Logger на основе конфигурации
func New(cfg *config.Config) *Logger {
	level := slog.LevelInfo
	switch cfg.LogLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	return &Logger{log: slog.New(handler)}
}

// Info логирует информационное сообщение
func (l *Logger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

// Error логирует ошибку
func (l *Logger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}

// Warn логирует предупреждение
func (l *Logger) Warn(msg string, args ...any) {
	l.log.Warn(msg, args...)
}

// With создаёт новый Logger с дополнительными полями
func (l *Logger) With(args ...any) *Logger {
	return &Logger{log: l.log.With(args...)}
}
