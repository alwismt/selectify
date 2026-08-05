package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/lmittmann/tint"
)

type formattedError struct {
	original error
	message  string
}

func (e *formattedError) Error() string {
	return e.message
}

func (e *formattedError) Unwrap() error {
	return e.original
}

var l = slog.Default()

func Init() {
	handler := tint.NewHandler(os.Stdout, &tint.Options{
		Level:      slog.LevelInfo,
		TimeFormat: "2006-01-02 15:04:05",
	})

	l = slog.New(handler)
	slog.SetDefault(l)
}

func Info(ctx context.Context, msg string) {
	_, _ = logWithCaller(ctx, slog.LevelInfo, "%s", msg)
}

func Warn(ctx context.Context, msg string) {
	_, _ = logWithCaller(ctx, slog.LevelWarn, "%s", msg)
}

func Error(ctx context.Context, err error, msg string) error {
	fullMsg := msg
	if err != nil {
		fullMsg = fmt.Sprintf("%s, %s", err.Error(), msg)
	}
	file, fn := logWithCaller(ctx, slog.LevelError, "%s", fullMsg)
	formattedMsg := fmt.Sprintf("%s.%s: %s", file, fn, fullMsg)
	return &formattedError{original: err, message: formattedMsg}
}

func Fatal(ctx context.Context, err error, msg string) {
	fullMsg := msg
	if err != nil {
		fullMsg = fmt.Sprintf("%s, %s", err.Error(), msg)
	}
	_, _ = logWithCaller(ctx, slog.LevelError, "%s", fullMsg)
	os.Exit(1)
}

func Infof(ctx context.Context, format string, args ...any) {
	_, _ = logWithCaller(ctx, slog.LevelInfo, format, args...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	_, _ = logWithCaller(ctx, slog.LevelWarn, format, args...)
}

func Errorf(ctx context.Context, err error, format string, args ...any) error {
	msg := fmt.Sprintf(format, args...)
	fullMsg := msg
	if err != nil {
		fullMsg = fmt.Sprintf("%s, %s", err.Error(), msg)
	}
	file, fn := logWithCaller(ctx, slog.LevelError, "%s", fullMsg)
	formattedMsg := fmt.Sprintf("%s.%s: %s", file, fn, fullMsg)
	return &formattedError{original: err, message: formattedMsg}
}

func Fatalf(ctx context.Context, err error, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	fullMsg := msg
	if err != nil {
		fullMsg = fmt.Sprintf("%s, %s", err.Error(), msg)
	}
	_, _ = logWithCaller(ctx, slog.LevelError, "%s", fullMsg)
	os.Exit(1)
}

func logWithCaller(ctx context.Context, level slog.Level, format string, args ...any) (string, string) {
	msg := fmt.Sprintf(format, args...)
	file, fn := callerInfo(3)
	l.Log(ctx, level, fmt.Sprintf("%s.%s: %s", file, fn, msg))
	return file, fn
}

func callerInfo(skip int) (string, string) {
	pc, file, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", "unknown"
	}
	fn := runtime.FuncForPC(pc)
	name := "unknown"
	if fn != nil {
		full := fn.Name()
		if idx := strings.LastIndex(full, "/"); idx >= 0 && idx+1 < len(full) {
			full = full[idx+1:]
		}
		if idx := strings.LastIndex(full, "."); idx >= 0 && idx+1 < len(full) {
			name = full[idx+1:]
		} else {
			name = full
		}
	}
	baseFile := path.Base(file)
	// Remove .go extension if present
	baseFile = strings.TrimSuffix(baseFile, ".go")
	// Remove _test suffix if present (for test files)
	baseFile = strings.TrimSuffix(baseFile, "_test")
	return baseFile, name
}
