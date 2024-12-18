package log

import (
	"context"
	"log/slog"
)

type NoopLogger struct{}

func (l NoopLogger) T(ctx context.Context, msg string, a ...any) {}
func (l NoopLogger) D(ctx context.Context, msg string, a ...any) {}
func (l NoopLogger) I(ctx context.Context, msg string, a ...any) {}
func (l NoopLogger) W(ctx context.Context, msg string, a ...any) {}
func (l NoopLogger) E(ctx context.Context, msg string, a ...any) {}
func (l NoopLogger) Slogger() *slog.Logger {
	return nil
}
