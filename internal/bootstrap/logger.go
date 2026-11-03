package bootstrap

import (
	"fmt"
	"io"
	"log/slog"
	"photo-web/internal/config"
	"strings"
)

func NewLogger(cfg *config.Config, output io.Writer) (*slog.Logger, error) {
	level, err := ParseLevel(cfg.Logger.Level)
	if err != nil {
		return nil, err
	}
	opt := &slog.HandlerOptions{Level: level, AddSource: cfg.Logger.AddSource}

	var handler slog.Handler
	format := strings.ToLower(cfg.Logger.Format)
	switch format {
	case "", "text":
		handler = slog.NewTextHandler(output, opt)
	case "json":
		handler = slog.NewJSONHandler(output, opt)
	default:
		return nil, fmt.Errorf("invalid logger format: %s", cfg.Logger.Format)
	}

	logger := slog.New(handler).With(slog.String("env", cfg.Server.Mode))
	return logger, nil
}

func ParseLevel(level string) (slog.Level, error) {
	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return slog.LevelError, fmt.Errorf("invalid log level: %s", level)
	}
}
