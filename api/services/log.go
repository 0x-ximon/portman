package services

import (
	"context"
	"io"
	"log/slog"
	"os"
)

type LoggerKey struct{}

func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerKey{}).(*slog.Logger); ok {
		return logger
	}

	return slog.Default()
}

func GetLogWriter() io.Writer {
	if os.Getenv("ENV") == "prod" {
		return os.Stdout
	}

	file, err := os.OpenFile("tmp/dev.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return os.Stdout
	}

	return file
}
