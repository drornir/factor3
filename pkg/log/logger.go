package log

import (
	"context"
	"log/slog"
)

type Logger interface {
	T(ctx context.Context, msg string, a ...any)
	D(ctx context.Context, msg string, a ...any)
	I(ctx context.Context, msg string, a ...any)
	W(ctx context.Context, msg string, a ...any)
	E(ctx context.Context, msg string, a ...any)
	Slogger() *slog.Logger
}
